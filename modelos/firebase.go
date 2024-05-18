package modelos

type TokenFirebase struct {
	Token     string
	IdUsuario int
}

func (TokenFirebase) TableName() string {
	return "token_firebase"
}
