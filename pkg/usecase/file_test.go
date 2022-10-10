package usecase

//func TestFileUseCase_Store(t *testing.T) {
//	mockRepo := new(mocks.FileRepository)
//	mockDS := domain.File{
//		Name: "1",
//		Path: "1",
//	}
//	t.Run("success", func(t *testing.T) {
//		tempMockDS := mockDS
//		tempMockDS.ID = 0
//		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.File{}, nil).Once()
//		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.File")).Return(int64(1), nil).Once()
//
//		u := NewFileUseCase(mockRepo, time.Second*3)
//		_, err := u.Store(context.TODO(), &tempMockDS)
//		assert.NoError(t, err)
//		assert.Equal(t, mockDS.Name, tempMockDS.Name)
//		mockRepo.AssertExpectations(t)
//	})
//
//	t.Run("name already exists", func(t *testing.T) {
//		tempMockDS := mockDS
//		tempMockDS.ID = 0
//		errWant := errors.New("name already exists")
//		var affectWant int64 = 0
//		mockRepo.On("GetByName", mock.Anything, mock.AnythingOfType("string")).Return(domain.File{Name: "1"}, nil).Once()
//		mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.File")).Return(affectWant, errWant).Once()
//
//		u := NewFileUseCase(mockRepo, time.Second*3)
//		affect, err := u.Store(context.TODO(), &tempMockDS)
//		assert.Error(t, err, errWant)
//		assert.Equal(t, affect, affectWant)
//	})
//}
//
//func TestFileUseCase_Update(t *testing.T) {
//	mockRepo := new(mocks.FileRepository)
//	mockDS := domain.File{
//		ID:   1,
//		Name: "1",
//		Path: "1",
//	}
//	t.Run("success", func(t *testing.T) {
//		want := mockDS
//		want.ID = 1
//		want.Name = "123456"
//		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.File")).Return(domain.File{}, nil).Once()
//		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.File")).Return(int64(1), nil).Once()
//		u := NewFileUseCase(mockRepo, time.Second*3)
//		affect, err := u.Update(context.TODO(), &want)
//		assert.NoError(t, err)
//		assert.Equal(t, affect, int64(1))
//		mockRepo.AssertExpectations(t)
//	})
//
//	t.Run("fail", func(t *testing.T) {
//		tempMock := mockDS
//		tempMock.ID = 1
//		wantErr := errors.New("name is not exists")
//		mockRepo.On("CheckExist", mock.Anything, mock.AnythingOfType("*domain.File")).Return(mockDS, nil).Once()
//		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.File")).Return(int64(0), wantErr).Once()
//		u := NewFileUseCase(mockRepo, time.Second*3)
//		affect, err := u.Update(context.TODO(), &tempMock)
//		assert.Error(t, err, wantErr)
//		assert.Equal(t, affect, int64(0))
//	})
//}
