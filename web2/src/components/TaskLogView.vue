<template>
  <div class="task-log-view" :class="{'task-log-view--with-message': item.message}">
    <v-alert
        type="info"
        text
        v-if="item.message"
    >{{ item.message }}
    </v-alert>

    <v-container class="pa-0 mb-2">
      <v-row no-gutters>
        <v-col>
          <v-list two-line subheader class="pa-0">
            <v-list-item class="pa-0">
              <v-list-item-content>
                <div>
                  <TaskStatus :status="item.status"/>
                </div>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-col>
        <v-col>
          <v-list two-line subheader class="pa-0">
            <v-list-item class="pa-0">
              <v-list-item-content>
                <v-list-item-title>Author</v-list-item-title>
                <v-list-item-subtitle>{{ user.name }}</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-col>
        <v-col>
          <v-list two-line subheader class="pa-0">
            <v-list-item class="pa-0">
              <v-list-item-content>
                <v-list-item-title>Started</v-list-item-title>
                <v-list-item-subtitle>
                  {{ item.start | formatDate }}
                </v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-col>
        <v-col>
          <v-list-item class="pa-0">
            <v-list-item-content>
              <v-list-item-title>Duration</v-list-item-title>
              <v-list-item-subtitle>
                {{ [item.start, item.end] | formatMilliseconds }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-col>
      </v-row>
    </v-container>
    <div class="task-log-records" ref="output">
      <div class="task-log-records__record" v-for="record in output" :key="record.id">
        <div class="task-log-records__time">
          {{ record.time | formatTime }}
        </div>
        <div class="task-log-records__output">{{ record.output }}</div>
      </div>
    </div>

    <v-btn
        color="error"
        style="position: absolute; bottom: 10px; right: 10px;"
        v-if="item.status === 'running' || item.status === 'waiting'"
        @click="stopTask()"
    >
      Stop
    </v-btn>
  </div>
</template>

<style lang="scss">

@import '~vuetify/src/styles/settings/_variables';

.task-log-view {
}

.task-log-records {
  background: black;
  color: white;
  height: calc(100vh - 250px);
  overflow: auto;
  font-family: monospace;
  margin: 0 -24px;
  padding: 5px 10px;
}

.task-log-view--with-message .task-log-records {
  height: calc(100vh - 300px);
}

.task-log-records__record {
  display: flex;
  flex-direction: row;
  justify-content: left;
}

.task-log-records__time {
  width: 120px;
  min-width: 120px;
}

.task-log-records__output {
  width: 100%;
}

@media #{map-get($display-breakpoints, 'sm-and-down')} {
  .task-log-records {
    height: calc(100vh - 340px);
  }

  .task-log-view--with-message .task-log-records {
    height: calc(100vh - 370px);
  }
}
</style>
<script>
import axios from 'axios';
import TaskStatus from '@/components/TaskStatus.vue';
import socket from '@/socket';

export default {
  components: { TaskStatus },
  props: {
    itemId: Number,
    projectId: Number,
  },
  data() {
    return {
      item: {},
      output: [],
      user: {},
    };
  },
  watch: {
    async itemId() {
      this.reset();
      await this.loadData();
    },

    async projectId() {
      this.reset();
      await this.loadData();
    },
  },
  async created() {
    socket.addListener((data) => this.onWebsocketDataReceived(data));
    await this.loadData();
  },

  methods: {
    async stopTask() {
      await axios({
        method: 'post',
        url: `/api/project/${this.projectId}/tasks/${this.itemId}/stop`,
        responseType: 'json',
      });
    },

    reset() {
      this.item = {};
      this.output = [];
      this.user = {};
    },

    onWebsocketDataReceived(data) {
      if (data.project_id !== this.projectId || data.task_id !== this.itemId) {
        return;
      }

      switch (data.type) {
        case 'update':
          Object.assign(this.item, {
            ...data,
            type: undefined,
          });
          break;
        case 'log':
          this.output.push(data);
          setTimeout(() => {
            this.$refs.output.scrollTop = this.$refs.output.scrollHeight;
          }, 200);
          break;
        default:
          break;
      }
    },

    async loadData() {
      this.item = (await axios({
        method: 'get',
        url: `/api/project/${this.projectId}/tasks/${this.itemId}`,
        responseType: 'json',
      })).data;

      this.output = (await axios({
        method: 'get',
        url: `/api/project/${this.projectId}/tasks/${this.itemId}/output`,
        responseType: 'json',
      })).data;

      this.user = (await axios({
        method: 'get',
        url: `/api/users/${this.item.user_id}`,
        responseType: 'json',
      })).data;
    },
  },
};
</script>
