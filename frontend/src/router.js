import { createRouter, createWebHistory } from 'vue-router'
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import ProblemsList from './views/ProblemsList.vue'
import ProblemDetail from './views/ProblemDetail.vue'
import AdminProblems from './views/AdminProblems.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  { path: '/problems', name: 'Problems', component: ProblemsList, meta: { requiresAuth: true } },
  { path: '/problems/:id', name: 'ProblemDetail', component: ProblemDetail, meta: { requiresAuth: true } },
  { path: '/admin/problems', name: 'AdminProblems', component: AdminProblems, meta: { requiresAuth: true, requiresAdmin: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !localStorage.getItem('token')) {
    return { name: 'Login' }
  }
  if (to.meta.requiresAdmin && localStorage.getItem('role') !== 'admin') {
    return { name: 'Problems' }
  }
})

export default router
