package services

import (
	"errors"

	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/repositories"
	"github.com/grahagandangr/nexboard-be/utils"
)

type BoardService struct {
	boardRepo     *repositories.BoardRepository
	workspaceRepo *repositories.WorkspaceRepository
	userRepo      *repositories.UserRepository
}

func NewBoardService(boardRepo *repositories.BoardRepository, workspaceRepo *repositories.WorkspaceRepository, userRepo *repositories.UserRepository) *BoardService {
	return &BoardService{
		boardRepo:     boardRepo,
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

// CreateBoard creates a new board in a workspace
func (s *BoardService) CreateBoard(userExternalID, workspaceExternalID string, req *models.BoardRequest) (*models.BoardResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	// Verify user is a member of the workspace
	_, err = s.workspaceRepo.GetMemberRole(w.ID, user.ID)
	if err != nil {
		return nil, errors.New("unauthorized: not a member of this workspace")
	}

	board := &models.Board{
		ExternalID:  utils.GenerateUUID(),
		WorkspaceID: w.ID,
		CreatedByID: &user.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.boardRepo.CreateBoard(board); err != nil {
		return nil, err
	}

	return &models.BoardResponse{
		ExternalID:          board.ExternalID,
		WorkspaceExternalID: w.ExternalID,
		Name:                board.Name,
		Description:         board.Description,
		CreatedAt:           board.CreatedAt,
		ModifiedAt:          board.ModifiedAt,
	}, nil
}

// GetWorkspaceBoards lists all boards in a workspace
func (s *BoardService) GetWorkspaceBoards(userExternalID, workspaceExternalID string) ([]*models.BoardResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	// Verify user is a member of the workspace
	_, err = s.workspaceRepo.GetMemberRole(w.ID, user.ID)
	if err != nil {
		return nil, errors.New("unauthorized: not a member of this workspace")
	}

	boards, err := s.boardRepo.GetBoardsByWorkspaceID(w.ID)
	if err != nil {
		return nil, err
	}

	var response []*models.BoardResponse
	for _, b := range boards {
		response = append(response, &models.BoardResponse{
			ExternalID:          b.ExternalID,
			WorkspaceExternalID: w.ExternalID,
			Name:                b.Name,
			Description:         b.Description,
			CreatedAt:           b.CreatedAt,
			ModifiedAt:          b.ModifiedAt,
		})
	}

	return response, nil
}

// GetBoard gets detailed info of a board
func (s *BoardService) GetBoard(userExternalID, boardExternalID string) (*models.BoardResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	b, err := s.boardRepo.GetBoardByExternalID(boardExternalID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	// Make sure user is member of workspace where this board resides
	_, err = s.workspaceRepo.GetMemberRole(b.WorkspaceID, user.ID)
	if err != nil {
		return nil, errors.New("unauthorized: not a member of the workspace for this board")
	}

	return &models.BoardResponse{
		ExternalID:          b.ExternalID,
		WorkspaceExternalID: b.WorkspaceExternalID,
		Name:                b.Name,
		Description:         b.Description,
		CreatedAt:           b.CreatedAt,
		ModifiedAt:          b.ModifiedAt,
	}, nil
}

// UpdateBoard updates a board
func (s *BoardService) UpdateBoard(userExternalID, boardExternalID string, req *models.BoardRequest) (*models.BoardResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	b, err := s.boardRepo.GetBoardByExternalID(boardExternalID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	// Verify user is a member of the workspace
	_, err = s.workspaceRepo.GetMemberRole(b.WorkspaceID, user.ID)
	if err != nil {
		return nil, errors.New("unauthorized: not a member of the workspace")
	}

	b.Name = req.Name
	b.Description = req.Description

	if err := s.boardRepo.UpdateBoard(b); err != nil {
		return nil, err
	}

	return &models.BoardResponse{
		ExternalID:          b.ExternalID,
		WorkspaceExternalID: b.WorkspaceExternalID,
		Name:                b.Name,
		Description:         b.Description,
		CreatedAt:           b.CreatedAt,
		ModifiedAt:          b.ModifiedAt,
	}, nil
}

// DeleteBoard deletes a board (and cascades its tasks)
func (s *BoardService) DeleteBoard(userExternalID, boardExternalID string) error {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return errors.New("user not found")
	}

	b, err := s.boardRepo.GetBoardByExternalID(boardExternalID)
	if err != nil {
		return errors.New("board not found")
	}

	// All workspace members can delete boards currently as per PRD (or maybe just owner/admin).
	// PRD says: "All workspace members can view boards and tasks". It doesn't restrict Board deletion inherently.
	// But logically, we check if they are at least a member.
	_, err = s.workspaceRepo.GetMemberRole(b.WorkspaceID, user.ID)
	if err != nil {
		return errors.New("unauthorized: not a member of the workspace")
	}

	return s.boardRepo.DeleteBoard(b.ID)
}
