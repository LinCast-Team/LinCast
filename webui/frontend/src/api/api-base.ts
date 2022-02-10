/* eslint-disable lines-between-class-members, @typescript-eslint/no-unused-vars */

export default abstract class APIBase {
  protected static readonly API_VERSION = 0;
  protected static readonly BASE_PATH = `/api/v${this.API_VERSION}`;
}
