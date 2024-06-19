<template>
    <div class="container">
        <div class="carousel-container">
            <n-carousel autoplay class="carousel" :dot-type="'line'" style="height: 180px">
                <img class="carousel-img" src="https://image-assets.mihuashi.com/artwork-carousel/5.jpg">
                <img class="carousel-img" src="https://image-assets.mihuashi.com/artwork-carousel/2.jpg">
                <img class="carousel-img" src="https://image-assets.mihuashi.com/artwork-carousel/3.jpg">
                <img class="carousel-img" src="https://image-assets.mihuashi.com/artwork-carousel/4.jpg">
            </n-carousel>
        </div>

    <div class="post-container">
        <n-button @click="showPostModal" class="post-button" size="large" strong secondary circle type="primary">
            发布投票
        </n-button>

        <n-modal v-model:show="postModalFlag" :style="{ width: '600px' }" preset="card" size="huge"
            :bordered="false">
            <n-form ref="formRef" :model="post_poll" label-placement="left" label-width="auto"
                require-mark-placement="right-hanging">
                <n-form-item label="投票主题" path="inputValue">
                    <n-input v-model:value="post_poll.description" placeholder="投票主题" />
                </n-form-item>
                
                <div style="display: flex; justify-content: flex-end">
                    <n-button round type="primary" @click="postSubmit">
                        确认
                    </n-button>
                </div>
            </n-form>
        </n-modal>
    </div>

    <div class="content-container">
        <n-grid :x-gap="12" :y-gap="8" :cols="4">
            <n-gi v-for="poll in polls" :key="poll.pollId">
                <CardVue :title="poll.description" :creator="poll.creator" @click="openPollModal(poll.pollId)"></CardVue>
            </n-gi>
        </n-grid>
    </div>

    <n-modal v-model:show="voteModalFlag" :style="{ width: '600px' }" preset="card" size="huge" :bordered="false">
        <h2>投票: {{ currentPoll.description }}</h2>
        <div v-for="candidate in currentPoll.candidates" :key="candidate">
            <n-button @click="vote(currentPoll.id, candidate)">{{ candidate }}</n-button>
        </div>
        <div style="display: flex; justify-content: flex-end; margin-top: 20px;">
            <n-button round type="primary" @click="joinPoll(currentPoll.id)">
                加入投票
            </n-button>
        </div>
    </n-modal>
</div>
</template>
<script setup>
import { ref, reactive } from 'vue';
import {
    NCarousel, NCarouselItem,
    NButton, NModal, NForm, NFormItem, NInput, NInputNumber,
    NGrid, NGi
} from 'naive-ui';
import CardVue from '@/components/Card.vue';
// import { BallotContract } from '@/contracts/Ballot';

// 初始化合约以及连接MetaMask
const contract = require('@truffle/contract');
const artifact = require('../assets/contracts/Ballot.json');
const Ballot = contract(artifact);
Ballot.setProvider(window.ethereum);

/*************************** polls相关 *******************************/
// 获得polls并在页面中显示
const polls = ref([]);
const loadPostedData = () => {
    Ballot.deployed().then((instance) => {
        return instance.getOngoingPolls.call();
    }).then((response) => {
        polls.value = response.map(poll => ({
            pollId: poll.pollId,
            description: poll.description,
            creator: poll.creator
        }));
    }).catch((err) => {
        console.log(err.message);
    });
};
loadPostedData();

/*************************** post相关 *******************************/
const postModalFlag = ref(false);
const post_poll = reactive({
    description: '',
    duration: 0
});

const showPostModal = () => {
    postModalFlag.value = true;
};

const postSubmit = () => {
    Ballot.deployed().then(async (instance) => {
        const accounts = await ethereum.request({ method: 'eth_requestAccounts' });
        const account = await accounts[0];
        await instance.createPoll(post_poll.description, { from: account });
    }).then(() => {
        loadPostedData();
        postModalFlag.value = false;
    }).catch((err) => {
        alert('error', err.message);
        console.log(err.message);
    });
};

/*************************** vote相关 *******************************/
const voteModalFlag = ref(false);
const currentPoll = reactive({
    id: 0,
    description: '',
    candidates: []
});

const openPollModal = (pollId) => {
    Ballot.deployed().then((instance) => {
        return instance.getCandidates.call(pollId);
    }).then((candidates) => {
        currentPoll.id = pollId;
        currentPoll.candidates = candidates;
        currentPoll.description = polls.value.find(p => p.pollId === pollId).description;
        voteModalFlag.value = true;
    }).catch((err) => {
        console.log(err.message);
    });
};

const vote = (pollId, candidate) => {
    Ballot.deployed().then(async (instance) => {
        const accounts = await ethereum.request({ method: 'eth_requestAccounts' });
        const account = await accounts[0];
        await instance.vote(pollId, candidate, { from: account });
    }).then(() => {
        voteModalFlag.value = false;
        loadPostedData();
    }).catch((err) => {
        alert('error', err.message);
        console.log(err.message);
    });
};

const joinPoll = (pollId) => {
    Ballot.deployed().then(async (instance) => {
        const accounts = await ethereum.request({ method: 'eth_requestAccounts' });
        const account = await accounts[0];
        await instance.joinPoll(pollId, { from: account });
    }).then(() => {
        voteModalFlag.value = false;
        loadPostedData();
    }).catch((err) => {
        alert('error', err.message);
        console.log(err.message);
    });
};
</script>
<style scoped>
.carousel-container {
    margin: 0 10%;
    padding: 1% 1%;
    background-color: honeydew;
    border-radius: 15px;
}

.carousel {
    width: 100%;
    border-radius: 15px;
}

.carousel-img {
    margin: 0 auto;
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.post-container {}

.post-button {
    position: fixed;
    padding: 30px;
    right: 40px;
    bottom: 40px;
    z-index: 1;
}

.modal {
    width: 600px;
}

.form {
    width: 600px;
}

.content-container {
    position: relative;
    min-height: 50vh;
    margin-top: 5vh;
    padding: 1% 5%;
    background-color: white;
    border-radius: 15px;
}
</style>