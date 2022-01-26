const state = {
  reqCancelMap: {},
  EsConnectID: 0,
  RefreshTab: '',
  LastSelectKey: [],
  activeName: '',
  ui: []
}

const mutations = {
  SET_EsConnect: (state, EsConnect) => {
    state.EsConnectID = EsConnect
  },
  SET_ReqCancelMap: (state, obj) => {
    state.reqCancelMap[obj['token']] = obj['fn']
  },
  DElETE_ReqCancelMap: (state, token) => {
    delete (state.reqCancelMap[token])
  },
  SET_RefreshTab: (state, RefreshTab) => {
    state.RefreshTab = RefreshTab
  },
  SET_LastSelectKey: (state, LastSelectKey) => {
    state.LastSelectKey = LastSelectKey
  },
  SET_ActiveName: (state, activeName) => {
    state.activeName = activeName
  },
  SET_Ui: (state, ui) => {
    state.ui = ui
  }
}

const actions = {
  SetEsConnect({ commit }, p) {
    commit('SET_EsConnect', p)
  },
  SET_ReqCancelMap({ commit }, p) {
    commit('SET_ReqCancelMap', p)
  },
  DElETE_ReqCancelMap({ commit }, p) {
    commit('DElETE_ReqCancelMap', p)
  },
  SETRefreshTab({ commit }, p) {
    commit('SET_RefreshTab', p)
  },
  SETActiveName({ commit }, p) {
    commit('SET_ActiveName', p)
  },
  SETLastSelectKey({ commit }, p) {
    commit('SET_LastSelectKey', p)
  },
  SETUI({ commit }, p) {
    commit('SET_Ui', p)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
