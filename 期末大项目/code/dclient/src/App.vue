<script setup>
import { useRouter } from 'vue-router';
import { h, ref, watch } from "vue";
import { NIcon, NMenu } from "naive-ui";
import {
  BookOutline as BookIcon,
  PersonOutline as PersonIcon,
  WineOutline as WineIcon
} from "@vicons/ionicons5";

const activeKey = ref("ballotsite")

function renderIcon(icon) {
  return () => h(NIcon, null, { default: () => h(icon) });
}
const menuOptions = [
  {
    label: "投票系统",
    key: "ballotsite",
  },
];

const router = useRouter();
watch(() => activeKey.value, () => {
    router.replace({
        path: activeKey.value
    })
})

</script>

<template>
  <div class="nav-container">
    <n-menu v-model:value="activeKey" mode="horizontal" :options="menuOptions" />
  </div>
  <div class="content-container">
    <router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component"></component>
      </keep-alive>
    </router-view>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
}

#app {
  width: 100%;
  height: fit-content;
  min-height: 100vh;
  background-color: aliceblue;
  background-image:
                    radial-gradient(closest-side, transparent 98%,rgb(234, 241, 246) 99%),
                    radial-gradient(closest-side,transparent 98%,rgb(234, 241, 246) 99%);
  background-position: 0 0px,40px 40px;
  background-size: 80px 80px;
}

.nav-container {
  margin: 0 1%;
  display: flex;
  justify-content: center;
  background-color: white;
  box-shadow: 2px 2px 5px rgb(173, 172, 172);
}

.content-container {
  margin-top: 1%;
  padding: 0 5%;
}
</style>

