package ndsim

import (
	"github.com/yanyiwu/gosimhash"
)

var gSimHasher gosimhash.Simhasher

const dictFileName = "jieba.dict.utf8"
const modelFileName = "hmm_model.utf8"
const idfFileName = "idf.utf8"
const stopWordFileName = "stop_words.utf8"

const hashTopN = 10

func initSimHasher() {
	gSimHasher = gosimhash.New(
		gConfig.DictPath+"/"+dictFileName,
		gConfig.DictPath+"/"+modelFileName,
		gConfig.DictPath+"/"+idfFileName,
		gConfig.DictPath+"/"+stopWordFileName)

}
