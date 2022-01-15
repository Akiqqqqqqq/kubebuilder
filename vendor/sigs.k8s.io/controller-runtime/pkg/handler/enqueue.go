/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/internal/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var enqueueLog = logf.RuntimeLog.WithName("eventhandler").WithName("EnqueueRequestForObject")

type empty struct{}

var _ EventHandler = &EnqueueRequestForObject{}

// EnqueueRequestForObject enqueues a Request containing the Name and Namespace of the object that is the source of the Event.
// (e.g. the created / deleted / updated objects Name and Namespace).  handler.EnqueueRequestForObject is used by almost all
// Controllers that have associated Resources (e.g. CRDs) to reconcile the associated Resource.

// EnqueueRequestForObject 是一个包含了作为事件源的对象的 Name 和 Namespace 的入队列的 Request
//（例如，created/deleted/updated 对象的 Name 和 Namespace）
// handler.EnqueueRequestForObject 几乎被所有关联资源（如 CRD）的控制器使用，以协调关联的资源
type EnqueueRequestForObject struct{}

// Create implements EventHandler.
func (e *EnqueueRequestForObject) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		enqueueLog.Error(nil, "CreateEvent received with no metadata", "event", evt)
		return
	}
	// 添加一个 Request 对象到工作队列
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      evt.Object.GetName(),
		Namespace: evt.Object.GetNamespace(),
	}})
}

// Update implements EventHandler.
func (e *EnqueueRequestForObject) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	switch {
	case evt.ObjectNew != nil:
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      evt.ObjectNew.GetName(),
			Namespace: evt.ObjectNew.GetNamespace(),
		}})
	case evt.ObjectOld != nil:
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      evt.ObjectOld.GetName(),
			Namespace: evt.ObjectOld.GetNamespace(),
		}})
	default:
		enqueueLog.Error(nil, "UpdateEvent received with no metadata", "event", evt)
	}
}

// Delete implements EventHandler.
func (e *EnqueueRequestForObject) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		enqueueLog.Error(nil, "DeleteEvent received with no metadata", "event", evt)
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      evt.Object.GetName(),
		Namespace: evt.Object.GetNamespace(),
	}})
}

// Generic implements EventHandler.
func (e *EnqueueRequestForObject) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		enqueueLog.Error(nil, "GenericEvent received with no metadata", "event", evt)
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      evt.Object.GetName(),
		Namespace: evt.Object.GetNamespace(),
	}})
}

// 通过 EnqueueRequestForObject 的 Create/Update/Delete 实现可以看出我们放入到工作队列中的元素不是以前默认的元素唯一的 KEY，
// 而是经过封装的 reconcile.Request 对象，当然通过这个对象也可以很方便获取对象的唯一标识 KEY
