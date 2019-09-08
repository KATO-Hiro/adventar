<template>
  <div>
    <GlobalHeader />

    <template v-if="calendar">
      <header :style="{ backgroundColor: calendarColor() }">
        <div class="inner">
          <h2>{{ calendar.title }} Advent Calendar {{ calendar.year }}</h2>
          <div>登録数 {{ calendar.entries.length }}/25人</div>
          <div>
            作成者 <UserIcon :user="calendar.owner" size="22" />
            {{ calendar.owner.name }}
          </div>
          <nuxt-link class="editBtn" :to="`/calendars/${calendar.id}/edit`" v-if="isOwnCalendar(calendar)">
            <font-awesome-icon icon="edit"></font-awesome-icon> 編集
          </nuxt-link>
        </div>
      </header>
      <main>
        <div>
          <VueMarkdown class="description">{{ calendar.description }}</VueMarkdown>
          <CalendarTable
            :calendar="calendar"
            :currentUser="$store.state.user"
            :onCreateentry="handleCreateEntry"
            :onUpdateEntry="handleUpdateEntry"
            :onDeleteEntry="handleDeleteEntry"
          ></CalendarTable>
          <EntryList :calendar="calendar"></EntryList>
        </div>
      </main>
    </template>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import VueMarkdown from "vue-markdown";
import { getCalendar, createEntry, updateEntry, deleteEntry } from "~/lib/GrpcClient";
import * as RestClient from "~/lib/RestClient";
import { calendarColor } from "~/lib/utils/Colors";
import { Calendar } from "~/types/adventar";
import { getToken } from "~/lib/Auth";
import GlobalHeader from "~/components/GlobalHeader.vue";
import UserIcon from "~/components/UserIcon.vue";
import CalendarTable from "~/components/CalendarTable.vue";
import EntryList from "~/components/EntryList.vue";

@Component({
  components: { GlobalHeader, UserIcon, VueMarkdown, CalendarTable, EntryList }
})
export default class extends Vue {
  calendar: Calendar | null = null;

  async asyncData({ params }) {
    if (process.server) {
      const calendar = await RestClient.getCalendar(params.id);
      return { calendar };
    }
  }

  async mounted() {
    this.calendar = await getCalendar(Number(this.$route.params.id));
  }

  async refetchCalendar() {
    this.calendar = await getCalendar(this.calendar!.id);
  }

  async handleCreateEntry(day: number): Promise<void> {
    const token = await getToken();
    const calendarId = this.calendar!.id;
    await createEntry({ calendarId, day, token });
    await this.refetchCalendar();
  }

  async handleUpdateEntry(entryId: number, { comment, url }: { comment: string; url: string }): Promise<void> {
    const token = await getToken();
    await updateEntry({ entryId, comment, url, token });
    await this.refetchCalendar();
  }

  async handleDeleteEntry(entryId: number): Promise<void> {
    const token = await getToken();
    await deleteEntry({ entryId, token });
    this.calendar = await getCalendar(this.calendar!.id);
    await this.refetchCalendar();
  }

  calendarColor(): string {
    return calendarColor(this.calendar!.id);
  }

  isOwnCalendar(calendar: Calendar): boolean {
    if (!this.$store.state.user) return false;
    if (!calendar.owner) return false;
    return calendar.owner.id === this.$store.state.user.id;
  }

  isFutureEntry(_): boolean {
    // TODO
    return false;
  }
}
</script>

<style lang="scss" scoped>
main > div {
  padding: 0;
}

header {
  color: #fff;
}

header > .inner {
  max-width: $content-max-width;
  padding: 30px 12px;
  margin: 0 auto;
  position: relative;
}

header .userIcon {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  vertical-align: middle;
}

header h2 {
  margin: 0 0 20px 0;
  font-size: 20px;
  font-weight: bold;
}

header .editBtn {
  padding: 10px 20px;
  font-size: 12px;
  display: inline-block;
  border-radius: 5px;
  cursor: pointer;
  text-decoration: none;
  border: none;
  outline: none;
  background: rgba(255, 255, 255, 0.85);
  color: #333;
  position: absolute;
  right: 12px;
  bottom: 30px;
}

.description {
  padding: 5px 10px;
}
</style>
