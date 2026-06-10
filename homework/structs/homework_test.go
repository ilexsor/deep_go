package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

// GamePersonStats битовая маска для хранения всех характеристик персонажа в одном uint32
// [31-22] - Mana (10 бит, значения 0-1000)
// [21-12] - Health (10 бит, значения 0-1000)
// [11-8]  - Level (4 бита, значения 0-10)
// [7-4]   - Experience (4 бита, значения 0-10)
// [3-0]   - Strength (4 бита, значения 0-10)
type GamePersonStats uint32

// GamePersonSocialStats битовая маска для хранения всех социальных характеристик персонажа в одном uint8
// [7-4]   - Respect (4 бита, значения 0-10)
// [3]     - HasFamily (1 бит)
// [2]     - HasHouse (1 бит)
// [1-0]   - Type (2 бита, значения 0-2)
type GamePersonSocialStats uint8

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
	clearLast10Bits         = 0xFFC00000
	clearMiddle10Bits       = 0x003FF000
	clearStrengthBits       = 0x0000000F
	clearExperienceBits     = 0x000000F0
	clearLevelBits          = 0x00000F00
	clearFirst4Bits         = 0xF0
	overflowProtection      = 0x3FF
	bits4OverflowProtection = 0x0F
	withHouseBits           = 0x04
	withFamilyBits          = 0x08
	withTypesBits           = 0x03
	shift22Bits             = 22
	shift12Bits             = 12
	shift4Bits              = 4
	shift8Bits              = 8
)

type GamePerson struct {
	x           int32                 // 4 байта
	y           int32                 // 4 байта
	z           int32                 // 4 байта
	gold        uint32                // 4 байта
	stats       GamePersonStats       // 4 байта
	name        [42]byte              // 42 байта
	socialStats GamePersonSocialStats // 1 байт
	hasWeapon   bool                  // 1 байт
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, opt := range options {
		opt(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	length := 0
	for length < len(p.name) && p.name[length] != 0 {
		length++
	}
	return string(p.name[:length])
}

func WithName(name string) Option {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) Option {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) Option {
	return func(person *GamePerson) {
		// Mana: биты 22-31 (10 бит)
		// Очищаем последние 10 бит 0xFFC00000 в двоичном виде: 11111111110000000000000000000000 (биты 22-31)
		// Если в маске (в числе справа) стоит 1, то в исходном числе этот бит принудительно сбрасывается в 0.
		person.stats &^= clearLast10Bits
		// 0x3FF = 0000 0000 0000 0000 0000 0011 1111 1111 Оператор & (битовое И) работает по правилу:
		// результат равен 1 только тогда, когда оба бита равны 1. Если один из битов 0, на выходе будет 0.
		// mana&0x3FF это трафарет для защиты от переполнения
		person.stats |= GamePersonStats(mana&overflowProtection) << shift22Bits
	}
}

func WithHealth(health int) Option {
	return func(person *GamePerson) {
		// Health: биты 12-21 (10 бит)
		person.stats &^= clearMiddle10Bits
		person.stats |= GamePersonStats(health&overflowProtection) << shift12Bits
	}
}

func WithRespect(respect int) Option {
	return func(person *GamePerson) {
		// Respect: биты 4-7 (4 бита)
		person.socialStats &^= clearFirst4Bits
		person.socialStats |= GamePersonSocialStats(respect&bits4OverflowProtection) << shift4Bits
	}
}

func WithStrength(strength int) Option {
	return func(person *GamePerson) {
		// Strength: биты 0-3 (4 бита)
		person.stats &^= clearStrengthBits
		person.stats |= GamePersonStats(strength & bits4OverflowProtection)
	}
}

func WithExperience(experience int) Option {
	return func(person *GamePerson) {
		// Experience: биты 4-7 (4 бита)
		person.stats &^= clearExperienceBits
		person.stats |= GamePersonStats(experience&bits4OverflowProtection) << shift4Bits
	}
}

func WithLevel(level int) Option {
	return func(person *GamePerson) {
		// Level: биты 8-11 (4 бита)
		person.stats &^= clearLevelBits
		person.stats |= GamePersonStats(level&bits4OverflowProtection) << shift8Bits
	}
}

func WithHouse() Option {
	return func(person *GamePerson) {
		person.socialStats |= withHouseBits // бит 2 для HasHouse
	}
}

func WithGun() Option {
	return func(person *GamePerson) {
		person.hasWeapon = true
	}
}

func WithFamily() Option {
	return func(person *GamePerson) {
		person.socialStats |= withFamilyBits // бит 3 для HasFamily
	}
}

func WithType(personType int) Option {
	return func(person *GamePerson) {
		// Type: биты 0-1 (2 бита)
		person.socialStats &^= withTypesBits
		person.socialStats |= GamePersonSocialStats(personType & withTypesBits)
	}
}

func (p *GamePerson) X() int    { return int(p.x) }
func (p *GamePerson) Y() int    { return int(p.y) }
func (p *GamePerson) Z() int    { return int(p.z) }
func (p *GamePerson) Gold() int { return int(p.gold) }

// Обратная операция по чтению бит с 31 по 22
// Сдвигаем все биты с 31 по 22 в начало, остальные биты зануляются
// И для безопасности применяем маску 0x3FF на получившееся число
func (p *GamePerson) Mana() int   { return int((p.stats >> shift22Bits) & overflowProtection) }
func (p *GamePerson) Health() int { return int((p.stats >> shift12Bits) & overflowProtection) }
func (p *GamePerson) Respect() int {
	return int((p.socialStats >> shift4Bits) & bits4OverflowProtection)
}
func (p *GamePerson) Strength() int    { return int(p.stats & bits4OverflowProtection) }
func (p *GamePerson) Experience() int  { return int((p.stats >> shift4Bits) & bits4OverflowProtection) }
func (p *GamePerson) Level() int       { return int((p.stats >> shift8Bits) & bits4OverflowProtection) }
func (p *GamePerson) HasHouse() bool   { return p.socialStats&withHouseBits != 0 }
func (p *GamePerson) HasGun() bool     { return p.hasWeapon }
func (p *GamePerson) HasFamilty() bool { return p.socialStats&withFamilyBits != 0 }
func (p *GamePerson) Type() int        { return int(p.socialStats & withTypesBits) }

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
