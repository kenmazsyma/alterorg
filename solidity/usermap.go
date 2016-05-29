package solidity
import (
    "bytes"
    "github.com/ethereum/go-ethereum/accounts/abi"
)
var Abi_User abi.ABI
var Abi_UserMap abi.ABI
func Init_usermap() error{
    var v *bytes.Buffer
    var er error
    v=bytes.NewBufferString(`[{"constant":false,"inputs":[{"name":"n","type":"string"}],"name":"changeName","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"node","type":"bytes"}],"name":"appendIpfsNode","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"node","type":"bytes"}],"name":"isExistNode","outputs":[{"name":"","type":"bool"}],"type":"function"},{"inputs":[{"name":"node","type":"bytes"},{"name":"n","type":"string"}],"type":"constructor"}]`)
    Abi_User, er=abi.JSON(v)
    if er != nil {
        return er
    }
    v=bytes.NewBufferString(`[{"constant":false,"inputs":[{"name":"adrs","type":"address"}],"name":"getUser","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"getAddresses","outputs":[{"name":"","type":"address[]"}],"type":"function"},{"constant":false,"inputs":[{"name":"node","type":"bytes"},{"name":"n","type":"string"}],"name":"reg","outputs":[],"type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"adrs","type":"address"},{"indexed":false,"name":"cont","type":"address"},{"indexed":false,"name":"isNew","type":"bool"}],"name":"onReg","type":"event"}]`)
    Abi_UserMap, er=abi.JSON(v)
    if er != nil {
        return er
    }
    return nil
}
var Bin_User="0x606060405233600060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908302179055506040516107d23803806107d2833981016040528080518201919060200180518201919060200150505b8060026000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100a757805160ff19168380011785556100d8565b828001600101855582156100d8579182015b828111156100d75782518260005055916020019190600101906100b9565b5b50905061010391906100e5565b808211156100ff57600081815060009055506001016100e5565b5090565b5050600160005080548060010182818154818355818115116101b6578183600052602060002091820191016101b59190610138565b808211156101b1576000818150805460018160011615610100020316600290046000825580601f1061016a57506101a7565b601f0160209004906000526020600020908101906101a69190610188565b808211156101a25760008181506000905550600101610188565b5090565b5b5050600101610138565b5090565b5b5050509190906000526020600020900160005b8490919091509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061021857805160ff1916838001178555610249565b82800160010185558215610249579182015b8281111561024857825182600050559160200191906001019061022a565b5b5090506102749190610256565b808211156102705760008181506000905550600101610256565b5090565b5050505b505061054a806102886000396000f360606040526000357c0100000000000000000000000000000000000000000000000000000000900480635353a2d81461004f578063daa02e3d146100a5578063ed99cf2a146100fb5761004d565b005b6100a36004808035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091905050610165565b005b6100f96004808035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091905050610216565b005b61014f6004808035906020019082018035906020019191908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505090909190505061038c565b6040518082815260200191505060405180910390f35b8060026000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101b457805160ff19168380011785556101e5565b828001600101855582156101e5579182015b828111156101e45782518260005055916020019190600101906101c6565b5b50905061021091906101f2565b8082111561020c57600081815060009055506001016101f2565b5090565b50505b50565b600160005080548060010182818154818355818115116102c7578183600052602060002091820191016102c69190610249565b808211156102c2576000818150805460018160011615610100020316600290046000825580601f1061027b57506102b8565b601f0160209004906000526020600020908101906102b79190610299565b808211156102b35760008181506000905550600101610299565b5090565b5b5050600101610249565b5090565b5b5050509190906000526020600020900160005b8390919091509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061032957805160ff191683800117855561035a565b8280016001018555821561035a579182015b8281111561035957825182600050559160200191906001019061033b565b5b5090506103859190610367565b808211156103815760008181506000905550600101610367565b5090565b5050505b50565b60006000600090505b6001600050805490508110156103e8576103cc600160005082815481101561000257906000526020600020900160005b50846103f7565b156103da57600191506103f1565b5b8080600101915050610395565b600091506103f1565b50919050565b600060006020604051908101604052806000815260200150600085925084915081518380546001816001161561010002031660029004905014151561043f5760009350610541565b600090505b828054600181600116156101000203166002900490508110156105385781818151811015610002579060200101517f010000000000000000000000000000000000000000000000000000000000000090047f0100000000000000000000000000000000000000000000000000000000000000028382815460018160011615610100020316600290048110156100025790908154600116156104f45790600052602060002090602091828204019190065b9054901a7f01000000000000000000000000000000000000000000000000000000000000000214151561052a5760009350610541565b5b8080600101915050610444565b60019350610541565b5050509291505056"
var Bin_UserMap="0x6060604052610d7d806100126000396000f360606040526000357c0100000000000000000000000000000000000000000000000000000000900480636f77926b1461004f578063a39fac1214610091578063e201c3fe146100e85761004d565b005b6100656004808035906020019091905050610185565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61009e60048050506101e0565b60405180806020018281038252838181518152602001915080519060200190602002808383829060006004602084601f0104600f02600301f1509050019250505060405180910390f35b6101836004808035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091908035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091905050610272565b005b6000600060005060008373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690506101db565b919050565b6020604051908101604052806000815260200150600160005080548060200260200160405190810160405280929190818152602001828054801561026357602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168152602001906001019080831161022f575b5050505050905061026f565b90565b600082826040516107d2806105ab8339018080602001806020018381038352858181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f1680156102e45780820380516001836020036101000a031916815260200191505b508381038252848181518152602001915080519060200190808383829060006004602084601f0104600f02600301f150905090810190601f16801561033d5780820380516001836020036101000a031916815260200191505b50945050505050604051809103906000f0905080600060005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690830217905550600073ffffffffffffffffffffffffffffffffffffffff16600060005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156105305760016000508054806001018281815481835581811511610478578183600052602060002091820191016104779190610459565b808211156104735760008181506000905550600101610459565b5090565b5b5050509190906000526020600020900160005b33909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff02191690830217905550507f45f6c878ece5d38c33ba9fb5beb53cd303dbc6d94641b6b1dcd04e9e13c7a27e33826001604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a16105a5565b7f45f6c878ece5d38c33ba9fb5beb53cd303dbc6d94641b6b1dcd04e9e13c7a27e33826000604051808473ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a15b5b50505056606060405233600060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908302179055506040516107d23803806107d2833981016040528080518201919060200180518201919060200150505b8060026000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100a757805160ff19168380011785556100d8565b828001600101855582156100d8579182015b828111156100d75782518260005055916020019190600101906100b9565b5b50905061010391906100e5565b808211156100ff57600081815060009055506001016100e5565b5090565b5050600160005080548060010182818154818355818115116101b6578183600052602060002091820191016101b59190610138565b808211156101b1576000818150805460018160011615610100020316600290046000825580601f1061016a57506101a7565b601f0160209004906000526020600020908101906101a69190610188565b808211156101a25760008181506000905550600101610188565b5090565b5b5050600101610138565b5090565b5b5050509190906000526020600020900160005b8490919091509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061021857805160ff1916838001178555610249565b82800160010185558215610249579182015b8281111561024857825182600050559160200191906001019061022a565b5b5090506102749190610256565b808211156102705760008181506000905550600101610256565b5090565b5050505b505061054a806102886000396000f360606040526000357c0100000000000000000000000000000000000000000000000000000000900480635353a2d81461004f578063daa02e3d146100a5578063ed99cf2a146100fb5761004d565b005b6100a36004808035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091905050610165565b005b6100f96004808035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091905050610216565b005b61014f6004808035906020019082018035906020019191908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505090909190505061038c565b6040518082815260200191505060405180910390f35b8060026000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101b457805160ff19168380011785556101e5565b828001600101855582156101e5579182015b828111156101e45782518260005055916020019190600101906101c6565b5b50905061021091906101f2565b8082111561020c57600081815060009055506001016101f2565b5090565b50505b50565b600160005080548060010182818154818355818115116102c7578183600052602060002091820191016102c69190610249565b808211156102c2576000818150805460018160011615610100020316600290046000825580601f1061027b57506102b8565b601f0160209004906000526020600020908101906102b79190610299565b808211156102b35760008181506000905550600101610299565b5090565b5b5050600101610249565b5090565b5b5050509190906000526020600020900160005b8390919091509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061032957805160ff191683800117855561035a565b8280016001018555821561035a579182015b8281111561035957825182600050559160200191906001019061033b565b5b5090506103859190610367565b808211156103815760008181506000905550600101610367565b5090565b5050505b50565b60006000600090505b6001600050805490508110156103e8576103cc600160005082815481101561000257906000526020600020900160005b50846103f7565b156103da57600191506103f1565b5b8080600101915050610395565b600091506103f1565b50919050565b600060006020604051908101604052806000815260200150600085925084915081518380546001816001161561010002031660029004905014151561043f5760009350610541565b600090505b828054600181600116156101000203166002900490508110156105385781818151811015610002579060200101517f010000000000000000000000000000000000000000000000000000000000000090047f0100000000000000000000000000000000000000000000000000000000000000028382815460018160011615610100020316600290048110156100025790908154600116156104f45790600052602060002090602091828204019190065b9054901a7f01000000000000000000000000000000000000000000000000000000000000000214151561052a5760009350610541565b5b8080600101915050610444565b60019350610541565b5050509291505056"
