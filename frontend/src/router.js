import { createRouter, createWebHistory } from 'vue-router'
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import ProblemsList from './views/ProblemsList.vue'
import ProblemDetail from './views/ProblemDetail.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  { path: '/problems', name: 'Problems', component: ProblemsList },
  { path: '/problems/:id', name: 'ProblemDetail', component: ProblemDetail },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
