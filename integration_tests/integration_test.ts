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
})
