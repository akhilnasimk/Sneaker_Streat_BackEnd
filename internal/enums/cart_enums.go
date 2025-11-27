package enums

type CartOperation string

const (
	OpInc CartOperation = "inc"
	OpDec CartOperation = "dec"
)

func (o CartOperation) IsValid() bool {
	return o == OpInc || o == OpDec
}