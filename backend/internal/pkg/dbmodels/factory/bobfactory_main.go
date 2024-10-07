// Code generated by BobGen psql v0.28.1. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

type Factory struct {
	baseAuthMods AuthModSlice
	baseCarMods  CarModSlice
	baseUserMods UserModSlice
}

func New() *Factory {
	return &Factory{}
}

func (f *Factory) NewAuth(mods ...AuthMod) *AuthTemplate {
	o := &AuthTemplate{f: f}

	if f != nil {
		f.baseAuthMods.Apply(o)
	}

	AuthModSlice(mods).Apply(o)

	return o
}

func (f *Factory) NewCar(mods ...CarMod) *CarTemplate {
	o := &CarTemplate{f: f}

	if f != nil {
		f.baseCarMods.Apply(o)
	}

	CarModSlice(mods).Apply(o)

	return o
}

func (f *Factory) NewUser(mods ...UserMod) *UserTemplate {
	o := &UserTemplate{f: f}

	if f != nil {
		f.baseUserMods.Apply(o)
	}

	UserModSlice(mods).Apply(o)

	return o
}

func (f *Factory) ClearBaseAuthMods() {
	f.baseAuthMods = nil
}

func (f *Factory) AddBaseAuthMod(mods ...AuthMod) {
	f.baseAuthMods = append(f.baseAuthMods, mods...)
}

func (f *Factory) ClearBaseCarMods() {
	f.baseCarMods = nil
}

func (f *Factory) AddBaseCarMod(mods ...CarMod) {
	f.baseCarMods = append(f.baseCarMods, mods...)
}

func (f *Factory) ClearBaseUserMods() {
	f.baseUserMods = nil
}

func (f *Factory) AddBaseUserMod(mods ...UserMod) {
	f.baseUserMods = append(f.baseUserMods, mods...)
}
