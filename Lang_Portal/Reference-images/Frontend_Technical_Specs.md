# Frontend Technical Spec

# Pages

### Dashboard /dashboard

#### Purpose
The purpose of this page is to provide summary of learning 
and act as the default page when a user visits the web-app

#### Components
- Last Study Session
    shows last activity used
    shows when last activity used
    summarizes wrong vs correct from last activity
    has a link to the group


- Study Progress
    - total words study
	- across all study session show the total words studied out of
	all possible words in our database   
	-display a mastery progress eg. 0%
    
- Quick Stats
    - success rate eg. 80%
    - total study sessions eg. 4
    - total active groups eg. 3
    - study streak eg. 4 days
- Sturt Studying Button
    - goes to study activities page

#### Needed API Endpoints 

- GET /dashboard/last_study_session
- GET /dashboard/study_progress
- GET /dashboard/quick-stats

### Study Activities Index/study-activities

#### Purpose
The purpose of this page is to show a collection of study activities 
with a thumbnail and its name,
to either launch or view the study activity.

#### Components


- Study Activity Card
    - show a thumbnail for the study activity
    - the name of the study activity
    - the launch button to take us to the launch page
    - the view page to view information about past
    study sessions for this study activity



#### Needed API Endpoints
 -GET  /api/study_activities

#### Study Activity Show '/study_activities/:id'
