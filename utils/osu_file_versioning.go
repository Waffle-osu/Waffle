package utils

import "github.com/Waffle-osu/osu-parser/osu_parser"

func VersionOsuFile(file osu_parser.OsuFile) int64 {
	switch file.Version {
	case 2:
		//Those builds pre-date build numbers
		return 0
	case 3:
		//Breaks got added in b73
		hasBreaks := len(file.Events.Events) > 0

		if hasBreaks {
			return 73
		} else {
			return 51
		}
	case 4:
		//b160 added the ability to have multiple timing points
		//Adds Custom Samplesets and Samples
		return 160
	case 5:
		//base v5 added a 24ms offset force on earlier than v5 maps
		//base v5 got added in b235
		//b402 added inherited timing points
		//b568 added SVs
		hasInheritedTimingPoints := func() bool {
			for _, point := range file.TimingPoints.TimingPoints {
				if point.InheritedTimingPoint {
					return true
				}
			}

			return false
		}()

		hasSvs := func() bool {
			for _, point := range file.TimingPoints.TimingPoints {
				if point.InheritedTimingPoint && point.BeatLength < 0 {
					return true
				}
			}

			return false
		}()

		if hasInheritedTimingPoints {
			if hasSvs {
				return 568
			} else {
				return 402
			}
		} else {
			return 235
		}
	case 6:
		//Stacking fixes and animation speeds fixed for storyboard sprites
		return 972
	case 7:
		//math error on multipart bezier sliders fixed
		return 1218
	case 8:
		//Sliderticks per beat added
		//Drain rate fix
		return 1650
	case 9:
		//Bezier is now the default slider type
		//Spinner new combos aren't forced anymore
		return 1688
	case 10:
		//Fixes bezier sliders being 1/50th shorter than they should
		return 1822
	case 11:
		//Mania hold notes added
		return 20121221
	case 12:
		//Per note hitsounding
		//Per note custom hitsounds
		//Per note volume
		return 20130220
	case 13:
		//Decimal diff settings
		return 20140612
	case 14:
		//Per-node samplesets on ctb sliders
		return 20150113
	}

	return 0
}
