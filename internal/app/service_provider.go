package app

import (
	"context"
	"github.com/vdyakova/grpc/internal/api/note"
	"github.com/vdyakova/grpc/internal/client/db"
	"github.com/vdyakova/grpc/internal/client/db/pg"
	"github.com/vdyakova/grpc/internal/client/db/transaction"
	"github.com/vdyakova/grpc/internal/closer"
	"github.com/vdyakova/grpc/internal/config"
	"github.com/vdyakova/grpc/internal/repository"
	noteRepository "github.com/vdyakova/grpc/internal/repository/note"
	"github.com/vdyakova/grpc/internal/repository_log"
	logRepository "github.com/vdyakova/grpc/internal/repository_log/note"
	"github.com/vdyakova/grpc/internal/service"
	noteService "github.com/vdyakova/grpc/internal/service/note"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	noteRepository repository.NoteRepository

	noteService   service.NoteService
	logRepository repository_log.LogRepository
	noteImpl      *note.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if s.noteRepository == nil {
		s.noteRepository = noteRepository.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(
			s.NoteRepository(ctx),
			s.TxManager(ctx),
			s.LogRepository(ctx),
		)
	}

	return s.noteService
}

func (s *serviceProvider) NoteImpl(ctx context.Context) *note.Implementation {
	if s.noteImpl == nil {
		s.noteImpl = note.NewImplementation(s.NoteService(ctx))
	}

	return s.noteImpl
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository_log.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepositoryLog(s.DBClient(ctx))
	}
	return s.logRepository
}
