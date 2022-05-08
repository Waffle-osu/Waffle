package helpers

func CalculateGlobalAccuracyOsu(hit50 uint64, hit100 uint64, hit300 uint64, hitGeki uint64, hitKatu uint64, hitMiss uint64) float32 {
	totalHits := hit50 + hit100 + hit300 + hitMiss
	perfectHits := float64(totalHits * 300)
	actualHits := float64(hit50*50 + hit100*100 + hit300*300)
	accuracy := actualHits / perfectHits

	return float32(accuracy)
}

func CalculateGlobalAccuracyTaiko(hit50 uint64, hit100 uint64, hit300 uint64, hitGeki uint64, hitKatu uint64, hitMiss uint64) float32 {
	totalHits := hit100 + hit300 + hitMiss
	perfectHits := float64(totalHits * 300)
	actualHits := float64(hit100*150 + hit300*300)
	accuracy := actualHits / perfectHits

	return float32(accuracy)
}

func CalculateGlobalAccuracyCatch(hit50 uint64, hit100 uint64, hit300 uint64, hitGeki uint64, hitKatu uint64, hitMiss uint64) float32 {
	totalHits := float64(hit50 + hit100 + hit300 + hitMiss + hitKatu)
	actualHits := float64(hit50 + hit100 + hit300)
	accuracy := actualHits / totalHits

	return float32(accuracy)
}
