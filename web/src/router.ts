import { createRouter, createWebHistory } from 'vue-router'
import ModulesView from './views/ModulesView.vue'
import CategoriesView from './views/CategoriesView.vue'
import LessonsView from './views/LessonsView.vue'
import LessonView from './views/LessonView.vue'
import PaymentView from './views/PaymentView.vue'
import ProfileView from './views/ProfileView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'modules',
      component: ModulesView,
    },
    {
      path: '/module/:id',
      name: 'categories',
      component: CategoriesView,
      props: true,
    },
    {
      path: '/category/:id',
      name: 'lessons',
      component: LessonsView,
      props: true,
    },
    {
      path: '/lesson/:id',
      name: 'lesson',
      component: LessonView,
      props: true,
    },
    {
      path: '/payment',
      name: 'payment',
      component: PaymentView,
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
    },
  ],
})

export default router
