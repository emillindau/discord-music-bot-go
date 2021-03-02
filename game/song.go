package game

type Song struct {
	artist string;
	name string;
	url string;
	FilePath string;
}

func newSong(a string, n string, u string) *Song {
	return &Song{
		artist: a,
		name: n,
		url: u,
	}
}