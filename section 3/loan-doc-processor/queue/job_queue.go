package queue

import "loan-doc-processor/model"

var JobChannel = make(chan model.DocumentJob, 100)
