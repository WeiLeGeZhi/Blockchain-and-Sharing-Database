import { createRouter, createWebHashHistory } from 'vue-router'
// import IndexVue from '@/views/Index.vue'
import BallotSiteVue from '@/views/BallotSite.vue'
// import AuctionVue from '@/views/Auction.vue'

// 路由规则
const routes = [
    {
        path: '/',
        name: 'default',
        component: BallotSiteVue
    },
    {
        path: '/ballotsite',
        name: 'ballotsite',
        component: BallotSiteVue
    },
]

// 路由实例
export default createRouter({
    history: createWebHashHistory(),
    routes,
})
