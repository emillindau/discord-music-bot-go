package game

type Song struct {
	id string;
	artist string;
	name string;
	downloadUrl string;
}

func newSong(a string, n string, du string, id string) *Song {
	return &Song{
		artist: a,
		name: n,
		downloadUrl: du,
		id: id,
	}
}