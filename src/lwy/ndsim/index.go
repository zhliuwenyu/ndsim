package ndsim

import (
	"errors"
	"os"
)

//IndexControl define index global info
type IndexControl struct {
	//Path               string
	DocIDFrom                DocID
	DocIDEnd                 DocID
	ForwardIndexFile         *os.File
	ForwardContentFile       *os.File
	LastForwardContentOffset DiskOffset
	ReverseIndexFile         *os.File
	ReverseContentFile       *os.File
	ReverseIndexMap          map[HashSign]ReverseIndex
	MaxReverseContentOffset  DiskOffset
}

//GIndexControl define global IndexControl object
var GIndexControl = IndexControl{ReverseIndexMap: make(map[HashSign]ReverseIndex)}

func initIndexControl() error {
	{ //open data file
		dir, err := os.Lstat(gConfig.DataPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(gConfig.DataPath, 0664)
				if err != nil {
					gLog.Fatal(err)
					return err
				}
			} else {
				gLog.Fatal(err)
				return err
			}
		} else if !dir.IsDir() {
			gLog.Fatal(gConfig.DataPath, "is already exist and not path")
			return errors.New("data path is not dir")
		}
		GIndexControl.ForwardIndexFile, err = os.OpenFile(gConfig.DataPath+"/"+gConfig.ForwardIndexFileName,
			os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			gLog.Fatal(err)
			return err
		}
		GIndexControl.ForwardContentFile, err = os.OpenFile(gConfig.DataPath+"/"+gConfig.ForwardContentFileName,
			os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			gLog.Fatal(err)
			return err
		}
		GIndexControl.ReverseIndexFile, err = os.OpenFile(gConfig.DataPath+"/"+gConfig.ReverseIndexFileName,
			os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			gLog.Fatal(err)
			return err
		}
		GIndexControl.ReverseContentFile, err = os.OpenFile(gConfig.DataPath+"/"+gConfig.ReverseContentFileName,
			os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			gLog.Fatal(err)
			return err
		}
	}
	if err := loadReverseIndexToMap(); err != nil {
		return err
	}
	{ //init LastForwardContentOffset
		info, err := os.Stat(gConfig.DataPath + "/" + gConfig.ForwardContentFileName)
		if err != nil {
			gLog.Fatal(err)
			return err
		}

		GIndexControl.LastForwardContentOffset = DiskOffset(info.Size())

	}
	return nil
}
