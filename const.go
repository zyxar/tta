package tta

import (
	"errors"
)

const (
	maxDepth        = 3
	maxBPS          = (maxDepth * 8)
	minBPS          = 16
	maxNCH          = 6
	fifoBufferSize  = 5120
	pcmBufferLength = 5120

	// TTA audio format
	formatSimple    = 1
	formatEncrypted = 2

	// Version of TTA lib
	Version = "0.0.2"
)

var (
	// TTA_CODEC_STATUS
	errOpen         = errors.New("can't open file")
	errFormat       = errors.New("not compatible file format")
	errFile         = errors.New("file is corrupted")
	errRead         = errors.New("can't read from input file")
	errWrite        = errors.New("can't write to output file")
	errSeek         = errors.New("file seek error")
	errPassword     = errors.New("password protected file")
	errNotSupported = errors.New("unsupported architecture")
	// Partial io
	errPartialWritten = errors.New("partial written")
)

var bitMask = []uint32{
	0x00000000, 0x00000001, 0x00000003, 0x00000007,
	0x0000000f, 0x0000001f, 0x0000003f, 0x0000007f,
	0x000000ff, 0x000001ff, 0x000003ff, 0x000007ff,
	0x00000fff, 0x00001fff, 0x00003fff, 0x00007fff,
	0x0000ffff, 0x0001ffff, 0x0003ffff, 0x0007ffff,
	0x000fffff, 0x001fffff, 0x003fffff, 0x007fffff,
	0x00ffffff, 0x01ffffff, 0x03ffffff, 0x07ffffff,
	0x0fffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff,
	0xffffffff,
}

var bitShift = []uint32{
	0x00000001, 0x00000002, 0x00000004, 0x00000008,
	0x00000010, 0x00000020, 0x00000040, 0x00000080,
	0x00000100, 0x00000200, 0x00000400, 0x00000800,
	0x00001000, 0x00002000, 0x00004000, 0x00008000,
	0x00010000, 0x00020000, 0x00040000, 0x00080000,
	0x00100000, 0x00200000, 0x00400000, 0x00800000,
	0x01000000, 0x02000000, 0x04000000, 0x08000000,
	0x10000000, 0x20000000, 0x40000000, 0x80000000,
	0x80000000, 0x80000000, 0x80000000, 0x80000000,
	0x80000000, 0x80000000, 0x80000000, 0x80000000,
}

var shift16 = bitShift[4:]

var crc32Table = [256]uint32{
	0x00000000, 0x77073096, 0xee0e612c, 0x990951ba, 0x076dc419, 0x706af48f,
	0xe963a535, 0x9e6495a3, 0x0edb8832, 0x79dcb8a4, 0xe0d5e91e, 0x97d2d988,
	0x09b64c2b, 0x7eb17cbd, 0xe7b82d07, 0x90bf1d91, 0x1db71064, 0x6ab020f2,
	0xf3b97148, 0x84be41de, 0x1adad47d, 0x6ddde4eb, 0xf4d4b551, 0x83d385c7,
	0x136c9856, 0x646ba8c0, 0xfd62f97a, 0x8a65c9ec, 0x14015c4f, 0x63066cd9,
	0xfa0f3d63, 0x8d080df5, 0x3b6e20c8, 0x4c69105e, 0xd56041e4, 0xa2677172,
	0x3c03e4d1, 0x4b04d447, 0xd20d85fd, 0xa50ab56b, 0x35b5a8fa, 0x42b2986c,
	0xdbbbc9d6, 0xacbcf940, 0x32d86ce3, 0x45df5c75, 0xdcd60dcf, 0xabd13d59,
	0x26d930ac, 0x51de003a, 0xc8d75180, 0xbfd06116, 0x21b4f4b5, 0x56b3c423,
	0xcfba9599, 0xb8bda50f, 0x2802b89e, 0x5f058808, 0xc60cd9b2, 0xb10be924,
	0x2f6f7c87, 0x58684c11, 0xc1611dab, 0xb6662d3d, 0x76dc4190, 0x01db7106,
	0x98d220bc, 0xefd5102a, 0x71b18589, 0x06b6b51f, 0x9fbfe4a5, 0xe8b8d433,
	0x7807c9a2, 0x0f00f934, 0x9609a88e, 0xe10e9818, 0x7f6a0dbb, 0x086d3d2d,
	0x91646c97, 0xe6635c01, 0x6b6b51f4, 0x1c6c6162, 0x856530d8, 0xf262004e,
	0x6c0695ed, 0x1b01a57b, 0x8208f4c1, 0xf50fc457, 0x65b0d9c6, 0x12b7e950,
	0x8bbeb8ea, 0xfcb9887c, 0x62dd1ddf, 0x15da2d49, 0x8cd37cf3, 0xfbd44c65,
	0x4db26158, 0x3ab551ce, 0xa3bc0074, 0xd4bb30e2, 0x4adfa541, 0x3dd895d7,
	0xa4d1c46d, 0xd3d6f4fb, 0x4369e96a, 0x346ed9fc, 0xad678846, 0xda60b8d0,
	0x44042d73, 0x33031de5, 0xaa0a4c5f, 0xdd0d7cc9, 0x5005713c, 0x270241aa,
	0xbe0b1010, 0xc90c2086, 0x5768b525, 0x206f85b3, 0xb966d409, 0xce61e49f,
	0x5edef90e, 0x29d9c998, 0xb0d09822, 0xc7d7a8b4, 0x59b33d17, 0x2eb40d81,
	0xb7bd5c3b, 0xc0ba6cad, 0xedb88320, 0x9abfb3b6, 0x03b6e20c, 0x74b1d29a,
	0xead54739, 0x9dd277af, 0x04db2615, 0x73dc1683, 0xe3630b12, 0x94643b84,
	0x0d6d6a3e, 0x7a6a5aa8, 0xe40ecf0b, 0x9309ff9d, 0x0a00ae27, 0x7d079eb1,
	0xf00f9344, 0x8708a3d2, 0x1e01f268, 0x6906c2fe, 0xf762575d, 0x806567cb,
	0x196c3671, 0x6e6b06e7, 0xfed41b76, 0x89d32be0, 0x10da7a5a, 0x67dd4acc,
	0xf9b9df6f, 0x8ebeeff9, 0x17b7be43, 0x60b08ed5, 0xd6d6a3e8, 0xa1d1937e,
	0x38d8c2c4, 0x4fdff252, 0xd1bb67f1, 0xa6bc5767, 0x3fb506dd, 0x48b2364b,
	0xd80d2bda, 0xaf0a1b4c, 0x36034af6, 0x41047a60, 0xdf60efc3, 0xa867df55,
	0x316e8eef, 0x4669be79, 0xcb61b38c, 0xbc66831a, 0x256fd2a0, 0x5268e236,
	0xcc0c7795, 0xbb0b4703, 0x220216b9, 0x5505262f, 0xc5ba3bbe, 0xb2bd0b28,
	0x2bb45a92, 0x5cb36a04, 0xc2d7ffa7, 0xb5d0cf31, 0x2cd99e8b, 0x5bdeae1d,
	0x9b64c2b0, 0xec63f226, 0x756aa39c, 0x026d930a, 0x9c0906a9, 0xeb0e363f,
	0x72076785, 0x05005713, 0x95bf4a82, 0xe2b87a14, 0x7bb12bae, 0x0cb61b38,
	0x92d28e9b, 0xe5d5be0d, 0x7cdcefb7, 0x0bdbdf21, 0x86d3d2d4, 0xf1d4e242,
	0x68ddb3f8, 0x1fda836e, 0x81be16cd, 0xf6b9265b, 0x6fb077e1, 0x18b74777,
	0x88085ae6, 0xff0f6a70, 0x66063bca, 0x11010b5c, 0x8f659eff, 0xf862ae69,
	0x616bffd3, 0x166ccf45, 0xa00ae278, 0xd70dd2ee, 0x4e048354, 0x3903b3c2,
	0xa7672661, 0xd06016f7, 0x4969474d, 0x3e6e77db, 0xaed16a4a, 0xd9d65adc,
	0x40df0b66, 0x37d83bf0, 0xa9bcae53, 0xdebb9ec5, 0x47b2cf7f, 0x30b5ffe9,
	0xbdbdf21c, 0xcabac28a, 0x53b39330, 0x24b4a3a6, 0xbad03605, 0xcdd70693,
	0x54de5729, 0x23d967bf, 0xb3667a2e, 0xc4614ab8, 0x5d681b02, 0x2a6f2b94,
	0xb40bbe37, 0xc30c8ea1, 0x5a05df1b, 0x2d02ef8d,
}

var crc64TableLow = [256]uint32{
	0x00000000, 0xa9ea3693, 0x53d46d26, 0xfa3e5bb5, 0x0e42ecdf, 0xa7a8da4c,
	0x5d9681f9, 0xf47cb76a, 0x1c85d9be, 0xb56fef2d, 0x4f51b498, 0xe6bb820b,
	0x12c73561, 0xbb2d03f2, 0x41135847, 0xe8f96ed4, 0x90e185ef, 0x390bb37c,
	0xc335e8c9, 0x6adfde5a, 0x9ea36930, 0x37495fa3, 0xcd770416, 0x649d3285,
	0x8c645c51, 0x258e6ac2, 0xdfb03177, 0x765a07e4, 0x8226b08e, 0x2bcc861d,
	0xd1f2dda8, 0x7818eb3b, 0x21c30bde, 0x88293d4d, 0x721766f8, 0xdbfd506b,
	0x2f81e701, 0x866bd192, 0x7c558a27, 0xd5bfbcb4, 0x3d46d260, 0x94ace4f3,
	0x6e92bf46, 0xc77889d5, 0x33043ebf, 0x9aee082c, 0x60d05399, 0xc93a650a,
	0xb1228e31, 0x18c8b8a2, 0xe2f6e317, 0x4b1cd584, 0xbf6062ee, 0x168a547d,
	0xecb40fc8, 0x455e395b, 0xada7578f, 0x044d611c, 0xfe733aa9, 0x57990c3a,
	0xa3e5bb50, 0x0a0f8dc3, 0xf031d676, 0x59dbe0e5, 0xea6c212f, 0x438617bc,
	0xb9b84c09, 0x10527a9a, 0xe42ecdf0, 0x4dc4fb63, 0xb7faa0d6, 0x1e109645,
	0xf6e9f891, 0x5f03ce02, 0xa53d95b7, 0x0cd7a324, 0xf8ab144e, 0x514122dd,
	0xab7f7968, 0x02954ffb, 0x7a8da4c0, 0xd3679253, 0x2959c9e6, 0x80b3ff75,
	0x74cf481f, 0xdd257e8c, 0x271b2539, 0x8ef113aa, 0x66087d7e, 0xcfe24bed,
	0x35dc1058, 0x9c3626cb, 0x684a91a1, 0xc1a0a732, 0x3b9efc87, 0x9274ca14,
	0xcbaf2af1, 0x62451c62, 0x987b47d7, 0x31917144, 0xc5edc62e, 0x6c07f0bd,
	0x9639ab08, 0x3fd39d9b, 0xd72af34f, 0x7ec0c5dc, 0x84fe9e69, 0x2d14a8fa,
	0xd9681f90, 0x70822903, 0x8abc72b6, 0x23564425, 0x5b4eaf1e, 0xf2a4998d,
	0x089ac238, 0xa170f4ab, 0x550c43c1, 0xfce67552, 0x06d82ee7, 0xaf321874,
	0x47cb76a0, 0xee214033, 0x141f1b86, 0xbdf52d15, 0x49899a7f, 0xe063acec,
	0x1a5df759, 0xb3b7c1ca, 0x7d3274cd, 0xd4d8425e, 0x2ee619eb, 0x870c2f78,
	0x73709812, 0xda9aae81, 0x20a4f534, 0x894ec3a7, 0x61b7ad73, 0xc85d9be0,
	0x3263c055, 0x9b89f6c6, 0x6ff541ac, 0xc61f773f, 0x3c212c8a, 0x95cb1a19,
	0xedd3f122, 0x4439c7b1, 0xbe079c04, 0x17edaa97, 0xe3911dfd, 0x4a7b2b6e,
	0xb04570db, 0x19af4648, 0xf156289c, 0x58bc1e0f, 0xa28245ba, 0x0b687329,
	0xff14c443, 0x56fef2d0, 0xacc0a965, 0x052a9ff6, 0x5cf17f13, 0xf51b4980,
	0x0f251235, 0xa6cf24a6, 0x52b393cc, 0xfb59a55f, 0x0167feea, 0xa88dc879,
	0x4074a6ad, 0xe99e903e, 0x13a0cb8b, 0xba4afd18, 0x4e364a72, 0xe7dc7ce1,
	0x1de22754, 0xb40811c7, 0xcc10fafc, 0x65facc6f, 0x9fc497da, 0x362ea149,
	0xc2521623, 0x6bb820b0, 0x91867b05, 0x386c4d96, 0xd0952342, 0x797f15d1,
	0x83414e64, 0x2aab78f7, 0xded7cf9d, 0x773df90e, 0x8d03a2bb, 0x24e99428,
	0x975e55e2, 0x3eb46371, 0xc48a38c4, 0x6d600e57, 0x991cb93d, 0x30f68fae,
	0xcac8d41b, 0x6322e288, 0x8bdb8c5c, 0x2231bacf, 0xd80fe17a, 0x71e5d7e9,
	0x85996083, 0x2c735610, 0xd64d0da5, 0x7fa73b36, 0x07bfd00d, 0xae55e69e,
	0x546bbd2b, 0xfd818bb8, 0x09fd3cd2, 0xa0170a41, 0x5a2951f4, 0xf3c36767,
	0x1b3a09b3, 0xb2d03f20, 0x48ee6495, 0xe1045206, 0x1578e56c, 0xbc92d3ff,
	0x46ac884a, 0xef46bed9, 0xb69d5e3c, 0x1f7768af, 0xe549331a, 0x4ca30589,
	0xb8dfb2e3, 0x11358470, 0xeb0bdfc5, 0x42e1e956, 0xaa188782, 0x03f2b111,
	0xf9cceaa4, 0x5026dc37, 0xa45a6b5d, 0x0db05dce, 0xf78e067b, 0x5e6430e8,
	0x267cdbd3, 0x8f96ed40, 0x75a8b6f5, 0xdc428066, 0x283e370c, 0x81d4019f,
	0x7bea5a2a, 0xd2006cb9, 0x3af9026d, 0x931334fe, 0x692d6f4b, 0xc0c759d8,
	0x34bbeeb2, 0x9d51d821, 0x676f8394, 0xce85b507,
}

var crc64TableHigh = [256]uint32{
	0x00000000, 0x42f0e1eb, 0x85e1c3d7, 0xc711223c, 0x49336645, 0x0bc387ae,
	0xccd2a592, 0x8e224479, 0x9266cc8a, 0xd0962d61, 0x17870f5d, 0x5577eeb6,
	0xdb55aacf, 0x99a54b24, 0x5eb46918, 0x1c4488f3, 0x663d78ff, 0x24cd9914,
	0xe3dcbb28, 0xa12c5ac3, 0x2f0e1eba, 0x6dfeff51, 0xaaefdd6d, 0xe81f3c86,
	0xf45bb475, 0xb6ab559e, 0x71ba77a2, 0x334a9649, 0xbd68d230, 0xff9833db,
	0x388911e7, 0x7a79f00c, 0xcc7af1ff, 0x8e8a1014, 0x499b3228, 0x0b6bd3c3,
	0x854997ba, 0xc7b97651, 0x00a8546d, 0x4258b586, 0x5e1c3d75, 0x1cecdc9e,
	0xdbfdfea2, 0x990d1f49, 0x172f5b30, 0x55dfbadb, 0x92ce98e7, 0xd03e790c,
	0xaa478900, 0xe8b768eb, 0x2fa64ad7, 0x6d56ab3c, 0xe374ef45, 0xa1840eae,
	0x66952c92, 0x2465cd79, 0x3821458a, 0x7ad1a461, 0xbdc0865d, 0xff3067b6,
	0x711223cf, 0x33e2c224, 0xf4f3e018, 0xb60301f3, 0xda050215, 0x98f5e3fe,
	0x5fe4c1c2, 0x1d142029, 0x93366450, 0xd1c685bb, 0x16d7a787, 0x5427466c,
	0x4863ce9f, 0x0a932f74, 0xcd820d48, 0x8f72eca3, 0x0150a8da, 0x43a04931,
	0x84b16b0d, 0xc6418ae6, 0xbc387aea, 0xfec89b01, 0x39d9b93d, 0x7b2958d6,
	0xf50b1caf, 0xb7fbfd44, 0x70eadf78, 0x321a3e93, 0x2e5eb660, 0x6cae578b,
	0xabbf75b7, 0xe94f945c, 0x676dd025, 0x259d31ce, 0xe28c13f2, 0xa07cf219,
	0x167ff3ea, 0x548f1201, 0x939e303d, 0xd16ed1d6, 0x5f4c95af, 0x1dbc7444,
	0xdaad5678, 0x985db793, 0x84193f60, 0xc6e9de8b, 0x01f8fcb7, 0x43081d5c,
	0xcd2a5925, 0x8fdab8ce, 0x48cb9af2, 0x0a3b7b19, 0x70428b15, 0x32b26afe,
	0xf5a348c2, 0xb753a929, 0x3971ed50, 0x7b810cbb, 0xbc902e87, 0xfe60cf6c,
	0xe224479f, 0xa0d4a674, 0x67c58448, 0x253565a3, 0xab1721da, 0xe9e7c031,
	0x2ef6e20d, 0x6c0603e6, 0xf6fae5c0, 0xb40a042b, 0x731b2617, 0x31ebc7fc,
	0xbfc98385, 0xfd39626e, 0x3a284052, 0x78d8a1b9, 0x649c294a, 0x266cc8a1,
	0xe17dea9d, 0xa38d0b76, 0x2daf4f0f, 0x6f5faee4, 0xa84e8cd8, 0xeabe6d33,
	0x90c79d3f, 0xd2377cd4, 0x15265ee8, 0x57d6bf03, 0xd9f4fb7a, 0x9b041a91,
	0x5c1538ad, 0x1ee5d946, 0x02a151b5, 0x4051b05e, 0x87409262, 0xc5b07389,
	0x4b9237f0, 0x0962d61b, 0xce73f427, 0x8c8315cc, 0x3a80143f, 0x7870f5d4,
	0xbf61d7e8, 0xfd913603, 0x73b3727a, 0x31439391, 0xf652b1ad, 0xb4a25046,
	0xa8e6d8b5, 0xea16395e, 0x2d071b62, 0x6ff7fa89, 0xe1d5bef0, 0xa3255f1b,
	0x64347d27, 0x26c49ccc, 0x5cbd6cc0, 0x1e4d8d2b, 0xd95caf17, 0x9bac4efc,
	0x158e0a85, 0x577eeb6e, 0x906fc952, 0xd29f28b9, 0xcedba04a, 0x8c2b41a1,
	0x4b3a639d, 0x09ca8276, 0x87e8c60f, 0xc51827e4, 0x020905d8, 0x40f9e433,
	0x2cffe7d5, 0x6e0f063e, 0xa91e2402, 0xebeec5e9, 0x65cc8190, 0x273c607b,
	0xe02d4247, 0xa2dda3ac, 0xbe992b5f, 0xfc69cab4, 0x3b78e888, 0x79880963,
	0xf7aa4d1a, 0xb55aacf1, 0x724b8ecd, 0x30bb6f26, 0x4ac29f2a, 0x08327ec1,
	0xcf235cfd, 0x8dd3bd16, 0x03f1f96f, 0x41011884, 0x86103ab8, 0xc4e0db53,
	0xd8a453a0, 0x9a54b24b, 0x5d459077, 0x1fb5719c, 0x919735e5, 0xd367d40e,
	0x1476f632, 0x568617d9, 0xe085162a, 0xa275f7c1, 0x6564d5fd, 0x27943416,
	0xa9b6706f, 0xeb469184, 0x2c57b3b8, 0x6ea75253, 0x72e3daa0, 0x30133b4b,
	0xf7021977, 0xb5f2f89c, 0x3bd0bce5, 0x79205d0e, 0xbe317f32, 0xfcc19ed9,
	0x86b86ed5, 0xc4488f3e, 0x0359ad02, 0x41a94ce9, 0xcf8b0890, 0x8d7be97b,
	0x4a6acb47, 0x089a2aac, 0x14dea25f, 0x562e43b4, 0x913f6188, 0xd3cf8063,
	0x5dedc41a, 0x1f1d25f1, 0xd80c07cd, 0x9afce626,
}

var shifts = [3]uint32{10, 9, 10}
