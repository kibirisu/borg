package repository

import (
	"context"
	"database/sql"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	Add(context.Context, db.AddStatusParams) error
	GetById(context.Context, int) (db.Status, error)
	GetByURI(context.Context, string) (db.Status, error)
	GetShares(context.Context, int) ([]db.Status, error)
	GetLocalStatuses(context.Context) ([]db.GetLocalStatusesRow, error)
	GetByIdWithMetadata(context.Context, int) (db.GetStatusByIdWithMetadataRow, error)
	DeleteByID(context.Context, int32) error
	GetComments(context.Context, int) ([]db.Status, error)
	GetPostComments(context.Context, int) ([]db.Status, error)
	Update(context.Context, db.UpdateStatusParams) (db.Status, error)
	Delete(context.Context, int) error
}

type statusRepository struct {
	q *db.Queries
}

var _ StatusRepository = (*statusRepository)(nil)

// Create implements StatusRepository.
func (r *statusRepository) Create(
	ctx context.Context,
	status db.CreateStatusParams,
) (db.Status, error) {
	return r.q.CreateStatus(ctx, status)
}

// Add implements StatusRepository.
func (r *statusRepository) Add(ctx context.Context, status db.AddStatusParams) error {
	return r.q.AddStatus(ctx, status)
}

// GetById implements StatusRepository.
func (r *statusRepository) GetById(ctx context.Context, id int) (db.Status, error) {
	return r.q.GetStatusById(ctx, int32(id))
}

// GetByURI implements StatusRepository.
func (r *statusRepository) GetByURI(ctx context.Context, uri string) (db.Status, error) {
	return r.q.GetStatusByURI(ctx, uri)
}

// GetById implements StatusRepository.
func (r *statusRepository) GetByIdWithMetadata(
	ctx context.Context,
	id int,
) (db.GetStatusByIdWithMetadataRow, error) {
	return r.q.GetStatusByIdWithMetadata(ctx, int32(id))
}

// GetShares implements StatusRepository.
func (r *statusRepository) GetShares(ctx context.Context, id int) ([]db.Status, error) {
	return r.q.GetStatusShares(ctx, sql.NullInt32{Int32: int32(id), Valid: true})
}

// GetLocalStatuses implements StatusRepository.
func (r *statusRepository) GetLocalStatuses(ctx context.Context) ([]db.GetLocalStatusesRow, error) {
	return r.q.GetLocalStatuses(ctx)
}

// DeleteByURI implements StatusRepository.
func (r *statusRepository) DeleteByURI(context.Context, string) error {
	panic("unimplemented")
}
// GetComments implements StatusRepository.
func (r *statusRepository) GetComments(ctx context.Context, id int) ([]db.Status, error) {
    rows, err := r.q.GetStatusComments(ctx, sql.NullInt32{Int32: int32(id), Valid: true})
    if err != nil {
        return nil, err
    }
    
    comments := make([]db.Status, 0, len(rows))
    for _, row := range rows {
        comments = append(comments, row.Status)
    }
    
    return comments, nil
}

// GetPostComments implements StatusRepository.
func (r *statusRepository) GetPostComments(ctx context.Context, id int) ([]db.Status, error) {
    rows, err := r.q.GetStatusComments(ctx, sql.NullInt32{Int32: int32(id), Valid: true})
    if err != nil {
        return nil, err
    }
    
    comments := make([]db.Status, 0, len(rows))
    for _, row := range rows {
        comments = append(comments, row.Status)
    }
    
    return comments, nil
}
// Update implements StatusRepository.
func (r *statusRepository) Update(ctx context.Context, params db.UpdateStatusParams) (db.Status, error) {
    return r.q.UpdateStatus(ctx, params)
}
// Delete implements StatusRepository.:
func (r *statusRepository) Delete(ctx context.Context, id int) error {
    return r.q.DeleteStatus(ctx, int32(id))
}
