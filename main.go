package main

func main() {
	ReadSettings()
	Parse(
		new(worldCmd),
		new(applyCmd),
		new(noopCmd),
		new(syncableCmd),

		new(fetchCmd),
		new(findNonPaniniCmd),
		new(statusCmd),
		new(pathCmd),
	)
}
