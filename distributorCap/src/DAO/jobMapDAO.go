package DAO

var StoredJobMap = NewJobMap()

func GetJob(key int) Job {
	tmp := StoredJobMap.Get(key)
	return tmp
}

func PutJob(key int, jo Job) {
	StoredJobMap.Put(key, jo)
}

func RemoveJob(key int) {
	StoredJobMap.Delete(key)
}
