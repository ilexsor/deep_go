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
		person.stats &^= 0xFFC00000
		// 0x3FF = 0000 0000 0000 0000 0000 0011 1111 1111Оператор & (битовое И) работает по правилу:
		// результат равен 1 только тогда, когда оба бита равны 1. Если один из битов 0, на выходе будет 0.
		// mana&0x3FF это трафарет для защиты от переполнения
		person.stats |= GamePersonStats(mana&0x3FF) << 22
	}
}

func WithHealth(health int) Option {
	return func(person *GamePerson) {
		// Health: биты 12-21 (10 бит)
		person.stats &^= 0x003FF000
		person.stats |= GamePersonStats(health&0x3FF) << 12
	}
}

func WithRespect(respect int) Option {
	return func(person *GamePerson) {
		// Respect: биты 4-7 (4 бита)
		person.socialStats &^= 0xF0
		person.socialStats |= GamePersonSocialStats(respect&0x0F) << 4
	}
}

func WithStrength(strength int) Option {
	return func(person *GamePerson) {
		// Strength: биты 0-3 (4 бита)
		person.stats &^= 0x0000000F
		person.stats |= GamePersonStats(strength & 0x0F)
	}
}

func WithExperience(experience int) Option {
	return func(person *GamePerson) {
		// Experience: биты 4-7 (4 бита)
		person.stats &^= 0x000000F0
		person.stats |= GamePersonStats(experience&0x0F) << 4
	}
}

func WithLevel(level int) Option {
	return func(person *GamePerson) {
		// Level: биты 8-11 (4 бита)
		person.stats &^= 0x00000F00
		person.stats |= GamePersonStats(level&0x0F) << 8
	}
}

func WithHouse() Option {
	return func(person *GamePerson) {
		person.socialStats |= 0x04 // бит 2 для HasHouse
	}
}

func WithGun() Option {
	return func(person *GamePerson) {
		person.hasWeapon = true
	}
}

func WithFamily() Option {
	return func(person *GamePerson) {
		person.socialStats |= 0x08 // бит 3 для HasFamily
	}
}

func WithType(personType int) Option {
	return func(person *GamePerson) {
		// Type: биты 0-1 (2 бита)
		person.socialStats &^= 0x03
		person.socialStats |= GamePersonSocialStats(personType & 0x03)
	}
}

func (p *GamePerson) X() int    { return int(p.x) }
func (p *GamePerson) Y() int    { return int(p.y) }
func (p *GamePerson) Z() int    { return int(p.z) }
func (p *GamePerson) Gold() int { return int(p.gold) }

// Обратная операция по чтению бит с 31 по 22
// Сдвигаем все биты с 31 по 22 в начало, остальные биты зануляются
// И для безопасности применяем маску 0x3FF на получившееся число
func (p *GamePerson) Mana() int        { return int((p.stats >> 22) & 0x3FF) }
func (p *GamePerson) Health() int      { return int((p.stats >> 12) & 0x3FF) }
func (p *GamePerson) Respect() int     { return int((p.socialStats >> 4) & 0x0F) }
func (p *GamePerson) Strength() int    { return int(p.stats & 0x0F) }
func (p *GamePerson) Experience() int  { return int((p.stats >> 4) & 0x0F) }
func (p *GamePerson) Level() int       { return int((p.stats >> 8) & 0x0F) }
func (p *GamePerson) HasHouse() bool   { return p.socialStats&0x04 != 0 }
func (p *GamePerson) HasGun() bool     { return p.hasWeapon }
func (p *GamePerson) HasFamilty() bool { return p.socialStats&0x08 != 0 }
func (p *GamePerson) Type() int        { return int(p.socialStats & 0x03) }

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
