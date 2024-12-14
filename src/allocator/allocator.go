package allocator

import (
	"errors"
)

type Allocator struct {
	registers []*Register
}

func NewAllocator(names []string) (*Allocator) {
  registers := make([]*Register, len(names))
  for i, name := range names {
    registers[i] = &Register{
      name: name,
      isFree: true,
    }
  }
  return &Allocator{registers: registers}
}

func (allocator *Allocator) Allocate() (*Register, error) {
	for _, reg := range allocator.registers {
		if reg.isFree {
			reg.setRegisterStatus(false)
			return reg, nil
		}
	}
	return nil, errors.New("No free register")
}

func (allocator *Allocator) Release(register *Register) error {
	for _, reg := range allocator.registers {
		if reg == register {
			reg.setRegisterStatus(true)
      reg.SetIsLoaded(false)
			return nil
		}
	}
	return errors.New("Unable to deallocate register. Not managed by allocator")
}
