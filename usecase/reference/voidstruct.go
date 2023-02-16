package reference

type VoidStruct struct {
}

// RepositoryRequestは引数を必ず指定しなければならないが引数が必要ない場合がある
//
// そういうときの穴埋め用として使う。型はVoidStruct
func Void() VoidStruct {
	return VoidStruct{}
}
