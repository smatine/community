import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),

  actions: {
    authenticate: function() {
      let { identification, password } = this.getProperties('identification', 'password');

      this.get('session')
      .authenticate(
        'authenticator:oauth2',
        identification,
        password
      )
      .then(() => {
        this.transitionToRoute('protected.index');
      })
      .catch((reason) => {
        this.set('errorMessage', reason.error || reason);
      });
    }
  }
});