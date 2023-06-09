/*
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * @Author: zhaozhenfang
 * @Date: 2022-01-19 17:30:12
 * @LastEditTime: 2022-06-27 15:20:32
 */
import { constantRouterMap, routes } from '@/router'
import router from '@/router';
import { getPermission } from "@/request/user"
import { hasPermission } from "@/utils/auth";

function filterAsyncRouter(routers, menus) {
  routers.forEach((route) => {
    if (!hasPermission(menus, route)) {
      route.meta.hidden = true;
    }
    route.children && filterAsyncRouter(route.children, menus)
  })
  return routers
}

const permission = {
  state: {
    routers: [],
    routes: routes,
    notfound: [],
    menus: [],
    operations: [],
    activePanel: ''
  },
  mutations: {
    SET_ROUTERS: (state, routers) => {
      state.routers = [...routers, ...constantRouterMap];
    },
    SET_MENUS: (state, menus) => {
      state.menus = menus
    },
    SET_OPERATIONS: (state, operations) => {
      state.operations = operations;
    },
    SET_LOGROUTERS: (state, routers) => {
      state.logRouters = routers;
    },
    SET_ACtiVEPANEL: (state, panel) => {
      state.activePanel = panel;
    },
    SET_NOTFOUND: (state, routers) => {
      state.notfound = routers;
    }
  },
  actions: {
    GenerateRoutes({ commit, state }, panel) {
      return new Promise(resolve => {
        const menus = state.menus;
        menus.push(panel);
        commit("SET_MENUS", menus)
        let routers = filterAsyncRouter(JSON.parse(JSON.stringify(routes)), menus)
        commit('SET_ROUTERS', routers)

        resolve()
      })
    },
    getPermission({ commit }, roles) {
      let roleId = roles.split(',').map(Number)
      return getPermission({ roleId: roleId }).then(res => {
        return new Promise((resolve, reject) => {
          if (res.data.code === 200) {
            let data = res.data.data;
            let { menu, button } = data;
            button.push("default_all");
            commit("SET_MENUS", menu.split(','));
            commit("SET_OPERATIONS", button);
            resolve()
          } else {
            reject()
          }
        })
      })
    },
    SetMenus({ commit }, menus) {
      commit("SET_MENUS", menus)
    },
    SetActivePanel({ commit }, panel) {
      commit("SET_ACtiVEPANEL", panel)
    },
    addRoute({ commit, state }, route) {
      routes[1].children.push(route);
      let routers = filterAsyncRouter(JSON.parse(JSON.stringify(routes)), state.menus)
      commit('SET_ROUTERS', routers)
      console.log(routers)
    }
  },
  getters: {
    addRoutes: state => {
      const { hostRouters, notfound } = state;
      return [...hostRouters, ...notfound];
    },
    getMenus: state => {
      return state.menus
    },
    getOperations: state => {
      return state.operations
    },
    getPaths: state => {
      return state.routers[1].children.filter(item => {
        return item.meta != undefined;
      }).map(item => item)
    }
  }
}

export default permission