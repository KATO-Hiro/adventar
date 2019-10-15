import fetch from "node-fetch";
import { Calendar, Entry } from "~/types/adventar";

export async function getCalendar(id: number): Promise<Calendar> {
  const response = await fetch(`${process.env.API_BASE_URL}/v1/calendars?calendar_id=${id}`);
  if (!response.ok) {
    throw new Error(response.status);
  }
  const body = await response.json();
  const calendar = body.calendar;
  calendar.entries = body.entries;
  return calendar;
}

export async function listEntries(userId: number): Promise<Entry[]> {
  const response = await fetch(`${process.env.API_BASE_URL}/v1/entries?user_id=${userId}`);
  const body = await response.json();
  return body.entries;
}
