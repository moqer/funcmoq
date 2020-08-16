package funcmoq

// func TestSurveyCron_500OnOngoingSurveysError1(t *testing.T) {
// 	// b := NewRepoBuilder()
// 	// b.WithOngoingSurveys().Returning(errors.New("problem"))
// 	// c := NewController(b.Build(), nil)

// 	repo := NewImr()
// 	repo.ongoingSurveys.With().Returning(errors.New("problem"))
// 	c := NewController(repo, nil)

// 	rr := httptest.NewRecorder()
// 	req, err := http.NewRequest("POST", "", nil)
// 	assert.Nil(t, err)
// 	handler := http.HandlerFunc(c.SurveyCron)

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
// }

// func TestSurveyCron_SendOnSameDay1(t *testing.T) {
// 	sType := SurveyType("fake")
// 	survey := Survey{CloseOn: time.Date(2020, 8, 7, 0, 0, 0, 0, time.UTC), ID: uuid.New(), Type: sType}

// 	repo := NewImr()
// 	repo.ongoingSurveys.With().Returning(nil, []interface{}{survey}...)
// 	repo.updateStatus.With(survey.ID, Closed).Returning(nil)

// 	sb := NewSurveyBuilder()
// 	sb.WithSend(survey).Returning(nil)
// 	sender := sb.Build()

// 	factory := NewSurveyFactory()
// 	factory.Register(sType, func() SendReminder { return sender })

// 	c := NewController(repo, factory)

// 	rr := httptest.NewRecorder()
// 	req, err := http.NewRequest("POST", "", nil)
// 	assert.Nil(t, err)
// 	handler := http.HandlerFunc(c.SurveyCron)
// 	c.now = func() time.Time { return survey.CloseOn }

// 	handler.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	assert.Equal(t, 1, sender.sendCount)
// }

// func NewImr() *imr {
// 	return &imr{
// 		ongoingSurveys: *NewResRegistry(),
// 		updateStatus:   *NewResRegistry(),
// 	}
// }

// type imr struct {
// 	ongoingSurveys ResRegistry
// 	updateStatus   ResRegistry
// }

// func (r imr) GetOngoingSurveys() (ss []Survey, err error) {
// 	res := r.ongoingSurveys.Get()
// 	// if res.err != nil {
// 	// 	return nil, res.err
// 	// }

// 	tmp := reflect.ValueOf(ss).Kind()
// 	_ = tmp

// 	if err := res.Retrieve(err, ss); err != nil {
// 		return nil, errors.New("can't retrieve")
// 	}

// 	// ss = res.values

// 	return ToSurvey(res.values)
// }

// func ToSurvey(a []interface{}) ([]Survey, error) {
// 	if len(a) == 0 {
// 		return nil, nil
// 	}
// 	surveys := make([]Survey, 0)
// 	for _, s := range a {
// 		tmp, ok := s.(Survey)
// 		if !ok {
// 			return nil, errors.New("The value can't be converted to a survey")
// 		}
// 		surveys = append(surveys, tmp)
// 	}
// 	return surveys, nil
// }

// func (r imr) DecreaseReminderCount(surveyID uuid.UUID) error {
// 	return nil
// }

// func (r imr) UpdateStatus(surveyID uuid.UUID, status SurveyStatus) error {
// 	res := r.updateStatus.Get(surveyID, status)
// 	return res.err
// }

// func (r imr) GetSurvey(ID uuid.UUID) (Survey, error) {
// 	return Survey{}, nil
// }

// //to be extracted in servicebp

// func NewResRegistry() *ResRegistry {
// 	return &ResRegistry{
// 		results: make(map[uint64]*Result),
// 	}
// }

// type ResRegistry struct {
// 	single  *Result
// 	results map[uint64]*Result
// }

// func (h ResRegistry) Get(key ...interface{}) *Result {
// 	hash, err := hashstructure.Hash(key, nil)
// 	if err != nil {
// 		panic(err) //testcode i think it's ok
// 	}
// 	result, exists := h.results[hash]
// 	if !exists {
// 		return &Result{
// 			err: errors.New("This key wasn't registered"),
// 		}
// 	}
// 	return result
// }

// func (h ResRegistry) With(key ...interface{}) *Result {
// 	var br Result
// 	hash, err := hashstructure.Hash(key, nil)
// 	if err != nil {
// 		panic(err) //testcode i think it's ok
// 	}
// 	h.results[hash] = &br
// 	return &br
// }

// type Result struct {
// 	err    error
// 	values []interface{}
// 	tmp    []byte
// }

// func (r *Result) Returning(err error, args ...interface{}) {
// 	if err != nil {
// 		r.err = err
// 	}
// 	if args != nil {
// 		v := reflect.ValueOf(args).Kind()
// 		x := reflect.ValueOf(args[0]).Kind()
// 		_, _ = x, v
// 		r.values = args
// 	}
// }

// func (r *Result) Retrieve(err error, args ...interface{}) error {
// 	err = r.err
// 	if len(r.values) != len(args) {
// 		return errors.New("cant convert object")
// 	}

// 	if err := json.Unmarshal(r.tmp, &args); err != nil {
// 		panic(err)
// 	}

// 	for _, val := range r.values {
// 		v := reflect.ValueOf(val)
// 		x := v.Kind()
// 		_ = x
// 		if v.Kind() == reflect.Slice {
// 			tmp := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
// 			reflect.Copy(tmp, v)
// 			v.Set(tmp)
// 		}
// 		// reflect.ValueOf(args[i]).is.Set(reflect.ValueOf(r.values[i]))
// 	}
// 	// args = r.values
// 	return nil
// }
