<template>
  <table>
    <thead>
      <tr>
        <th>SUN</th>
        <th>MON</th>
        <th>TUE</th>
        <th>WED</th>
        <th>THU</th>
        <th>FRI</th>
        <th>SAT</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(cells, i) in rows" :key="i">
        <td v-for="cell in cells" :key="cell.day">
          <div class="popup" v-if="cell.entryable && displayedPopupCellDay === cell.day" @click.stop>
            <span role="button" class="popupCloseBtn" @click="displayedPopupCellDay = null">
              <font-awesome-icon icon="times" />
            </span>
            <form @submit.prevent="handleSubmitPopupForm(cell.entry)">
              <div class="popupRow">
                <font-awesome-icon icon="comment" />
                <input
                  type="text"
                  placeholder="記事の内容の予定などを入力してください"
                  :value="cell.entry.comment"
                  :ref="`inputComment${cell.entry.day}`"
                />
              </div>
              <div class="popupRow">
                <font-awesome-icon icon="link" />
                <input
                  type="text"
                  placeholder="登録した日になったらURLを入力してください"
                  :value="cell.entry.url"
                  :ref="`inputUrl${cell.entry.day}`"
                />
              </div>
              <div class="popupAction">
                <button class="submit" type="submit">保存</button>
                <span role="button" class="cancel" @click="handleClickDeleteEntry(cell.entry)">
                  登録をキャンセル
                </span>
              </div>
            </form>
          </div>
          <div v-if="cell.entryable" class="inner">
            <span class="day">{{ cell.day }}</span>
            <div v-if="cell.entry">
              <div class="entryUser">
                <UserIcon :user="cell.entry.owner" size="50" />
                <div>{{ cell.entry.owner.name }}</div>
              </div>
              <span
                class="editBtn"
                role="button"
                v-if="isOwnEntry(cell.entry)"
                @click.stop="handleClickEditEntry(cell.entry)"
              >
                <font-awesome-icon icon="edit" />
              </span>
              <span
                class="cancelBtn"
                role="button"
                v-if="forceCancelable(cell.entry)"
                @click="handleClickDeleteEntry(cell.entry)"
              >
                <font-awesome-icon icon="times" />
              </span>
            </div>
            <div v-else class="entryForm">
              <button @click="handleClickCreateEntry(cell.day)">登録</button>
            </div>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import UserIcon from "~/components/UserIcon.vue";

@Component({
  components: { UserIcon }
})
export default class extends Vue {
  displayedPopupCellDay: number | null;

  mounted() {
    document.addEventListener("click", this.handleClickDocument);
  }

  destroyed() {
    document.removeEventListener("click", this.handleClickDocument);
  }

  handleClickDocument() {
    this.displayedPopupCellDay = null;
  }

  data() {
    return {
      displayedPopupCellDay: null
    };
  }

  forceCancelable(entry: Entry): boolean {
    return (
      this.calendar !== null &&
      !this.isOwnEntry(entry) &&
      this.isOwnCalendar(this.calendar) &&
      !entry.url &&
      // TODO: Fix to JST
      dayjs(new Date(this.calendar.year, 12, entry.day)).add(1, "day") < dayjs(new Date())
    );
  }

  isOwnCalendar(calendar: Calendar): boolean {
    if (!this.$store.state.user) return false;
    if (!calendar.owner) return false;
    return calendar.owner.id === this.$store.state.user.id;
  }

  isOwnEntry(entry: Entry): boolean {
    if (!this.$store.state.user) return false;
    if (!entry.owner) return false;
    return entry.owner.id === this.$store.state.user.id;
  }

  isFutureEntry(_): boolean {
    // TODO
    return false;
  }

  async handleClickCreateEntry(day): Promise<void> {
    const calendarId = this.calendar!.id;
    const token = await getToken();
    await createEntry({ calendarId, day, token });
    await this.refetchCalendar();
    this.displayedPopupCellDay = day;
  }

  async handleClickEditEntry(entry: Entry): Promise<void> {
    this.displayedPopupCellDay = entry.day || null;
  }

  async handleClickDeleteEntry(entry: Entry): Promise<void> {
    if (!window.confirm("登録をキャンセルします")) return;
    const token = await getToken();
    await deleteEntry({ entryId: entry.id, token });
    this.calendar = await getCalendar(this.calendar!.id);
    this.displayedPopupCellDay = null;
    await this.refetchCalendar();
  }

  async handleSubmitPopupForm(entry: Entry): Promise<void> {
    const entryId = entry.id;
    const comment = this.$refs[`inputComment${entry.day}`][0].value;
    const url = this.$refs[`inputUrl${entry.day}`][0].value;
    const token = await getToken();
    await updateEntry({ entryId, comment, url, token });
    await this.refetchCalendar();
    this.displayedPopupCellDay = null;
  }
}
</script>

<style lang="scss" scoped>
table {
  width: 100%;
  table-layout: fixed;
  border-collapse: collapse;
}

thead th {
  border: 1px solid #f1f1f1;
  border-bottom: 3px solid #f1f1f1;
  background: #aaa;
  height: 50px;
  color: #fff;
  text-align: left;
  padding-left: 15px;
  font-size: 15px;
  text-transform: uppercase;
}

thead th:first-child {
  background: #e7998e;
}

thead th:last-child {
  background: #87a3d0;
}

td {
  border: 1px solid #e5e5e5;
  background-color: #fff;
  vertical-align: top;
  position: relative;
}

.inner {
  position: relative;
  min-height: 140px;
}

.day {
  position: absolute;
  top: 12px;
  left: 15px;
  font-size: 24px;
  font-weight: bold;
  color: #aaa;
}

td:first-child .day {
  color: #e7998e;
}

td:last-child .day {
  color: #87a3d0;
}

.entryForm {
  padding-top: 60px;
  text-align: center;
}

.entryForm button {
  padding: 0 15px;
  height: 32px;
  line-height: 32px;
  font-size: 15px;
  display: inline-block;
  border-radius: 5px;
  margin: 0;
  font-family: inherit;
  cursor: pointer;
  text-decoration: none;
  border: none;
  outline: none;
  font-weight: bold;
  background-color: #e0e0e0;
  background-image: linear-gradient(to bottom, #e0e0e0, #e0e0e0 50%, #d3d3d3 50%, #d3d3d3);
  color: #999;
}

.entryForm button:hover {
  color: #fff;
  background-color: #ef7266;
  background-image: linear-gradient(to bottom, #ef7266, #ef7266 50%, #e45541 50%, #e45541);
}

.entryUser {
  text-align: center;
  position: relative;
  top: 50px;
  color: #666;
}

.entryUser .userIcon {
  width: 50px;
  height: 50px;
  border-radius: 50px;
}

.cancelBtn {
  position: absolute;
  top: 15px;
  right: 10px;
  color: #e5004f;
  cursor: pointer;
}

.editBtn {
  position: absolute;
  top: 15px;
  right: 10px;
  color: #13b5b1;
  cursor: pointer;
}
.popup {
  width: 350px;
  height: 130px;
  border: 1px solid #ccc;
  padding: 35px 20px 15px 20px;
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  color: #666;
  top: -152px;
  left: -126px;
  position: absolute;
  background-color: #fff;
  z-index: 1;
}

.popup:before,
.popup:after {
  top: 100%;
  left: 50%;
  border: solid transparent;
  content: " ";
  height: 0;
  width: 0;
  position: absolute;
  pointer-events: none;
}
.popup:before {
  border-color: rgba(204, 204, 204, 0);
  border-top-color: #cccccc;
  border-width: 12px;
  margin-left: -12px;
}
.popup:after {
  border-color: rgba(255, 255, 255, 0);
  border-top-color: #ffffff;
  border-width: 10px;
  margin-left: -10px;
}

.popup .popupCloseBtn {
  position: absolute;
  top: 10px;
  right: 30px;
  cursor: pointer;
}

.popupRow {
  padding: 5px 0;
}

.popup input[type="text"] {
  width: 310px;
  padding: 5px;
  border: 1px solid #e5e5e5;
  outline: none;
  font-size: 13px;
  border-radius: 3px;
  font-family: inherit;
}

.popup .popupAction {
  margin: 10px 0 0 20px;
}

.popup .submit {
  font-size: 14px;
  color: #fff;
  border: 2px solid #28a745;
  background-color: #28a745;
  padding: 6px 20px;
  border-radius: 10px;
  cursor: pointer;
  margin-right: 100px;
}

.popup .cancel {
  font-size: 14px;
  color: #dc3545;
  padding: 6px 20px;
  cursor: pointer;
}
</style>
