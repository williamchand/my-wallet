package usecase

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/williamchand/my-wallet/wallet"
	"github.com/williamchand/my-wallet/author"
	"github.com/williamchand/my-wallet/models"
)

type walletUsecase struct {
	walletRepo    wallet.Repository
	authorRepo     author.Repository
	contextTimeout time.Duration
}

// NewWalletUsecase will create new an walletUsecase object representation of wallet.Usecase interface
func NewWalletUsecase(a wallet.Repository, ar author.Repository, timeout time.Duration) wallet.Usecase {
	return &walletUsecase{
		walletRepo:    a,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *walletUsecase) fillAuthorDetails(c context.Context, data []*models.Wallet) ([]*models.Wallet, error) {

	g, ctx := errgroup.WithContext(c)

	// Get the author's id
	mapAuthors := map[int64]models.Author{}

	for _, wallet := range data {
		mapAuthors[wallet.Author.ID] = models.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan *models.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != nil {
			mapAuthors[author.ID] = *author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *walletUsecase) Fetch(c context.Context, cursor string, num int64) ([]*models.Wallet, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listWallet, nextCursor, err := a.walletRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	listWallet, err = a.fillAuthorDetails(ctx, listWallet)
	if err != nil {
		return nil, "", err
	}

	return listWallet, nextCursor, nil
}

func (a *walletUsecase) GetByID(c context.Context, id int64) (*models.Wallet, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.walletRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor
	return res, nil
}

func (a *walletUsecase) Update(c context.Context, ar *models.Wallet) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.walletRepo.Update(ctx, ar)
}

func (a *walletUsecase) GetByTitle(c context.Context, title string) (*models.Wallet, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.walletRepo.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return nil, err
	}
	res.Author = *resAuthor

	return res, nil
}

func (a *walletUsecase) Store(c context.Context, m *models.Wallet) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedWallet, _ := a.GetByTitle(ctx, m.Title)
	if existedWallet != nil {
		return models.ErrConflict
	}

	err := a.walletRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *walletUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedWallet, err := a.walletRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existedWallet == nil {
		return models.ErrNotFound
	}
	return a.walletRepo.Delete(ctx, id)
}
