export class ErrorResponse {
  constructor(
    public error: {
      type: string,
      message: string,
    }
  ){}
}
