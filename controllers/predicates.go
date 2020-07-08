package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var (
	podPredicates = predicate.Funcs{
		CreateFunc: func(createEvent event.CreateEvent) bool {
			return PodRequiresSecret(createEvent.Object)
		},
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			return PodRequiresSecret(deleteEvent.Object)
		},
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			return PodRequiresSecret(updateEvent.ObjectOld) ||
				PodRequiresSecret(updateEvent.ObjectNew)
		},
		GenericFunc: func(genericEvent event.GenericEvent) bool {
			return PodRequiresSecret(genericEvent.Object)
		},
	}

	secretPredicates = predicate.Funcs{
		CreateFunc: func(event event.CreateEvent) bool {
			return IsManagedSecret(event.Meta) || IsSource(event.Meta)
		},
		DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
			return IsManagedSecret(deleteEvent.Meta) || IsSource(deleteEvent.Meta)
		},
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			return (IsManagedSecret(updateEvent.MetaOld) || IsManagedSecret(updateEvent.MetaNew)) ||
				(IsSource(updateEvent.MetaOld) || IsSource(updateEvent.MetaNew))
		},
		GenericFunc: func(genericEvent event.GenericEvent) bool {
			return IsManagedSecret(genericEvent.Meta) || IsSource(genericEvent.Meta)
		},
	}
)
