/* eslint-disable lines-between-class-members, @typescript-eslint/no-unused-vars */

abstract class APIBase {
  protected readonly API_VERSION = 0;
  protected readonly BASE_PATH = `/api/v${this.API_VERSION}`;
}
