import {CounterService} from "./service.pb";
import {expect} from 'chai'

describe("test grpc-gateway-ts communication", () => {
  it("unary request", async () => {
    const result = await CounterService.Increment({counter: 199}, {pathPrefix: "http://localhost:8081"})

    expect(result.result).to.equal(200)
  })

  it('streaming request', async () => {
    const response = [] as number[]
    await CounterService.StreamingIncrements({counter: 1}, (resp) => response.push(resp.result), {pathPrefix: "http://localhost:8081"})

    expect(response).to.deep.equal([2,3,4,5,6])
  })

  it('http get check request', async () => {
    const result = await CounterService.HTTPGet({num: 10}, {pathPrefix: "http://localhost:8081"})
    expect(result.result).to.equal(11)
  })

  it('http post body check request with nested body path', async () => {
    const result = await CounterService.HTTPPostWithNestedBodyPath({a: 10, req: { b: 15 }}, {pathPrefix: "http://localhost:8081"})
    expect(result.result).to.equal(25)
  })


  it('http post body check request with star in path', async () => {
    const result = await CounterService.HTTPPostWithStarBodyPath({a: 10, req: { b: 15 }, c: 23}, {pathPrefix: "http://localhost:8081"})
    expect(result.result).to.equal(48)
  })
})
