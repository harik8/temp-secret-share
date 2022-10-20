import { createWebHistory, createRouter } from "vue-router"
import Read from "../Read.vue"
import Main from "../Main.vue"

const routes = [
  {
    path: "/",
    name: "Main",
    component: Main,
  },
  {
    path: "/secret/:secretid",
    name: "Secret",
    component: Read,
  },
]

const router = createRouter({
    history: createWebHistory(),
    routes
  })
  

export default router