package tui

type EditOperation struct {
	Type   string // "insert", "delete", "replace"
	Text   string
	Row    int
	Col    int
	Offset int
	Length int
}

type UndoRedoStack struct {
	undoStack []*EditOperation
	redoStack []*EditOperation
	maxSize   int
}

func NewUndoRedoStack(maxSize int) *UndoRedoStack {
	if maxSize <= 0 {
		maxSize = 500 // Default limit
	}
	return &UndoRedoStack{
		undoStack: make([]*EditOperation, 0, maxSize),
		redoStack: make([]*EditOperation, 0, maxSize),
		maxSize:   maxSize,
	}
}

func (u *UndoRedoStack) Push(op *EditOperation) {
	u.undoStack = append(u.undoStack, op)

	// Enforce max size
	if len(u.undoStack) > u.maxSize {
		u.undoStack = u.undoStack[len(u.undoStack)-u.maxSize:]
	}

	// Clear redo stack on new operation
	u.redoStack = make([]*EditOperation, 0, u.maxSize)
}

func (u *UndoRedoStack) Undo() *EditOperation {
	if len(u.undoStack) == 0 {
		return nil
	}

	op := u.undoStack[len(u.undoStack)-1]
	u.undoStack = u.undoStack[:len(u.undoStack)-1]

	u.redoStack = append(u.redoStack, op)
	if len(u.redoStack) > u.maxSize {
		u.redoStack = u.redoStack[len(u.redoStack)-u.maxSize:]
	}

	return op
}

func (u *UndoRedoStack) Redo() *EditOperation {
	if len(u.redoStack) == 0 {
		return nil
	}

	op := u.redoStack[len(u.redoStack)-1]
	u.redoStack = u.redoStack[:len(u.redoStack)-1]

	u.undoStack = append(u.undoStack, op)
	if len(u.undoStack) > u.maxSize {
		u.undoStack = u.undoStack[len(u.undoStack)-u.maxSize:]
	}

	return op
}

func (u *UndoRedoStack) CanUndo() bool {
	return len(u.undoStack) > 0
}

func (u *UndoRedoStack) CanRedo() bool {
	return len(u.redoStack) > 0
}

func (u *UndoRedoStack) Clear() {
	u.undoStack = make([]*EditOperation, 0, u.maxSize)
	u.redoStack = make([]*EditOperation, 0, u.maxSize)
}

func (u *UndoRedoStack) UndoCount() int {
	return len(u.undoStack)
}

func (u *UndoRedoStack) RedoCount() int {
	return len(u.redoStack)
}
