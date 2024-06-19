<template>
    <n-card hoverable embedded size="small">
        <div class="title-container">
            <div class="title">{{ title }}</div>
        </div>
    </n-card>
</template>

<script setup>
import { NCard } from 'naive-ui';

const props = defineProps({
    title: String
})
</script>

<style scoped>
.n-card {
    width: 100%;
}

.title-container {
    border-radius: 10px;
    position: relative;
    overflow: hidden;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 10px;
}

.title {
    width: fit-content;
    padding: 0 10%;
    background-color: rgba(255, 255, 255, 0.7);
    border-radius: 10px;
    font-size: large;
    color: black;
}
</style>
