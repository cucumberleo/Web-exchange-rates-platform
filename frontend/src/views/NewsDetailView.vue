<template>
  <el-container>
    <el-main>
      <el-card v-if="article" class="article-detail">
        <h1>{{ article.Title }}</h1>
        <p>{{ article.Content }}</p>
        <div>
          <el-button :type="isLiked ? 'info' : 'primary'" @click="likeArticle">
            {{ isLiked ? 'Cancel like' : 'Like' }}
          </el-button>
          <p>Likes: {{ likes }}</p>
        </div>
      </el-card>
      <div v-else class="no-data">You must register and login to see this page!</div>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted, onUpdated, onBeforeUpdate, onBeforeMount } from "vue";
import { useRoute } from "vue-router";
import axios from "../axios";
import type { Article, Like } from "../types/Article";

const article = ref<Article | null>(null);
const route = useRoute();
const likes = ref<number>(0)
const isLiked = ref<boolean>(false)
interface LikesResponse{
  likes: number;
}
interface isLikedResponse{
  is_liked: boolean
}
interface likeArticleResponse{
  message: string;
  likes: number;
  is_liked: boolean
}
const { id } = route.params;

const fetchArticle = async () => {
  try {
    const response = await axios.get<Article>(`/articles/${id}`);
    article.value = response.data;
  } catch (error) {
    console.error("Failed to load article:", error);
  }
};

const likeArticle = async () => {
  try {
    const res = await axios.post<Like>(`articles/${id}/like`)
    likes.value = res.data.likes
    isLiked.value = res.data.is_liked
  } catch (error) {
    console.log('Error Liking article:', error)
  }
};

const fetchLike = async ()=>{
  try{
    const res = await axios.get<Like>(`articles/${id}/like`)
    likes.value = res.data.likes
  }catch(error){
    console.log('Error fetching likes:', error)
  }
}
// 查看用户是否已经点赞，去设置初始按钮状态
const checkInitialLikeStatus = async ()=>{
  try{
    const res = await axios.get<isLikedResponse>(`articles/${id}/isliked`);
    isLiked.value = res.data.is_liked
  }catch(error){
    console.log('Error checking if liked:',error)
  }
}
onMounted(fetchArticle);
onMounted(fetchLike);
onMounted(checkInitialLikeStatus)
</script>

<style scoped>
.article-detail {
  margin: 20px 0;
}

.no-data {
  text-align: center;
  font-size: 1.2em;
  color: #999;
}
</style>
