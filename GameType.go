package GameFight

type X_GUID int

func (guid X_GUID) IsValid() bool {
	return int(guid) > 0
}


type X_INT int
type X_UINT uint
type X_FLOAT float64
