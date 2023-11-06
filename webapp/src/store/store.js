import {createStore} from 'vuex'

const store = createStore({
  state(){
    return {
      user:null
    }
  },
  mutation:{
    setUser(state,user){
      state.user = user
    }
  },
  actions:{
    login({commit},{username,password}){
      const userData = {username,password}
      commit('setUser',userData)
      localStorage.setItem('user',JSON.stringify(userData));
    },
    logout ({ commit }) {
      commit('setUser', null);
      localStorage.removeItem('user');
    },
    checkUserLogin ({ commit }) {
      const user = JSON.parse(localStorage.getItem('user'));
      if (user) {
        commit('setUser', user);
      }
    }
  }
})
export default store
