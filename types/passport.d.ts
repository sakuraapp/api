declare module 'passport/lib/http/request' {
    interface Request {
        [key: string]: () => void
    }

    const req: Request

    export default req
}
