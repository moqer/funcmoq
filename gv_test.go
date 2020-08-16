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
// 	repo.ongoingSurveys.With().Returning(nil, []Survey{survey, survey})
// 	repo.updateStatus.With(survey.ID, Closed).Returning(nil)

// 	sender := NewTestSurvey()

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
// 	if res.err != nil {
// 		return nil, res.err
// 	}

// 	ss = make([]Survey, 0)
// 	if err := res.Retrieve(err, &ss); err != nil {
// 		return nil, errors.New("can't retrieve")
// 	}

// 	return ss, err
// }

// func (r *Result) Retrieve(err error, args ...interface{}) error {
// 	err = r.err
// 	if len(r.values) != len(args) {
// 		return errors.New("cant convert object")
// 	}
// 	// var y int
// 	// if reflect.TypeOf(r.values).ConvertibleTo(reflect.TypeOf(args)) {
// 	// 	y = reflect.Copy(reflect.ValueOf(r.values), reflect.ValueOf(args))
// 	// }

// 	for i := range r.values {
// 		v := reflect.ValueOf(r.values[i])
// 		t := reflect.TypeOf(r.values[i])
// 		x := v.Kind()
// 		y := reflect.ValueOf(args[i]).Kind()
// 		log.Println(reflect.ValueOf(args[i]).CanSet())
// 		if t.AssignableTo(reflect.TypeOf(args[i])) {
// 			reflect.Copy(reflect.ValueOf(args[i]), v)
// 			// reflect.ValueOf(args[i]).Set(v)
// 		} else if t.ConvertibleTo(reflect.TypeOf(args[i])) {
// 			reflect.ValueOf(args[i]).Set(v.Convert(reflect.TypeOf(args[i])))
// 		}
// 		// if t.ConvertibleTo(reflect.TypeOf(args[i])) {
// 		// 	y = reflect.Copy(reflect.ValueOf(r.values[i]), reflect.ValueOf(args[i]))
// 		// }
// 		log.Println(x, y)
// 	}

// 	// for _, val := range r.values {
// 	// 	v := reflect.ValueOf(val)
// 	// 	x := v.Kind()
// 	// 	_ = x
// 	// 	if v.Kind() == reflect.Slice {
// 	// 		tmp := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
// 	// 		reflect.Copy(tmp, v)
// 	// 		v.Set(tmp)
// 	// 	}
// 	// 	// reflect.ValueOf(args[i]).is.Set(reflect.ValueOf(r.values[i]))
// 	// }
// 	// args = r.values
// 	return nil
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

//to be extracted in servicebp
