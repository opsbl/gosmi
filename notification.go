package gosmi

import "github.com/opsbl/gosmi/types"

type Notification struct {
	SmiNode
	Objects []SmiNode
}

func (n SmiNode) AsNotification() Notification {
	return Notification{
		SmiNode: n,
		Objects: n.GetNotificationObjects(),
	}
}

func (n SmiNode) GetNotificationObjects() (objects []SmiNode) {
	for element := n.instance.smiInst.GetFirstElement(n.smiNode); element != nil; element = n.instance.smiInst.GetNextElement(element) {
		object := n.instance.smiInst.GetElementNode(element)
		if object == nil {
			// TODO: error
			return
		}
		if object.NodeKind != types.NodeNotification {
			// TODO: error
		}
		objects = append(objects, CreateNode(n.instance, object))
	}
	return
}
