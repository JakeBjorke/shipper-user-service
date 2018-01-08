package main

import (
	"context"
	"log"

	pb "github.com/jakebjorke/shipper-user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo         Repository
	tokenService Authable
}

//Get is used to get a user
func (srv *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		return err
	}

	res.User = user
	return nil
}

//GetAll gets all of the current system users.
func (srv *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.repo.GetAll()
	if err != nil {
		return err
	}

	res.Users = users
	return nil
}

//Auth is used to authenticate a user.
func (srv *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.repo.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		return err
	}

	//Compares our given password against the hashed password from the DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(user)
	if err != nil {
		return err
	}

	res.Token = token
	return nil
}

//Create is used to create a new user
func (srv *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPass)
	if err := srv.repo.Create(req); err != nil {
		return err
	}

	res.User = req
	return nil
}

//ValidateToken validates a token
func (srv *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	//ignoring for now...
	return nil
}
