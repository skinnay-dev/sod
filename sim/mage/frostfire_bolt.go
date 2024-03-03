package mage

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (mage *Mage) registerFrostfireBoltSpell() {
	if !mage.HasRune(proto.MageRune_RuneBeltFrostfireBolt) {
		return
	}

	level := float64(mage.Level)
	baseCalc := 13.828124 + 0.018012*level + 0.044141*level*level
	baseDamageLow := baseCalc * 2.58
	baseDamageHigh := baseCalc * 3.0
	baseDotDamage := baseCalc * .08
	spellCoeff := .857
	castTime := time.Second * 3
	manaCost := .14

	numTicks := int32(3)
	tickLength := time.Second * 3

	mage.FrostfireBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:  core.ActionID{SpellID: int32(proto.MageRune_RuneBeltFrostfireBolt)},
		SpellCode: SpellCode_MageFrostfireBolt,
		// TODO: Multi-school spells
		SpellSchool:  core.SpellSchoolFrost, // | core.SpellSchoolFire
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | SpellFlagChillSpell | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Frostfire Bolt",
			},
			NumberOfTicks: numTicks,
			TickLength:    tickLength,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseDotDamage
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.MageCritMultiplier(0),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
